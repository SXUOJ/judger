package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/Sxu-Online-Judge/judge/config"
	"github.com/Sxu-Online-Judge/judge/pkg/rlimit"
	"github.com/Sxu-Online-Judge/judge/pkg/seccomp"
	"github.com/Sxu-Online-Judge/judge/runner"
	"github.com/Sxu-Online-Judge/judge/sandbox"
)

var (
	showPrint   = true
	showDetails = true
	// showDetails    = false

	unsafe         bool
	allowProc      = false
	workDir        string
	pType          = "default"
	inputFileName  = "input.txt"
	outputFileName = "output.txt"
	errorFileName  = "error.txt"
	timeLimit      = uint64(1)
	realTimeLimit  = uint64(1)
	memoryLimit    = uint64(128)
	outputLimit    = uint64(128)
	stackLimit     = uint64(128)
)

func main() {
	if realTimeLimit < timeLimit {
		realTimeLimit = timeLimit + 2
	}
	if stackLimit > memoryLimit {
		stackLimit = memoryLimit
	}
	if workDir == "" {
		workDir, _ = os.Getwd()
	}

	rt, err := start(os.Args[1:])
	if rt == nil {
		rt = &runner.Result{
			Status: runner.StatusSystemError,
		}
	}
	if err == nil && rt.Status != runner.StatusNormal {
		err = rt.Status
	}

}

func start(args []string) (*runner.Result, error) {
	var (
		err      error
		execFile uintptr
		rt       runner.Result
	)

	allow, trace, h := config.GetConf(pType, allowProc)

	files, err := prepareFiles(inputFileName, outputFileName, errorFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare files: %v", err)
	}
	defer closeFiles(files)

	// if not defined, then use the original value
	fds := make([]uintptr, len(files))
	for i, f := range files {
		if f != nil {
			fds[i] = f.Fd()
		} else {
			fds[i] = uintptr(i)
		}
	}

	rlimits := rlimit.RLimits{
		CPU:         timeLimit,
		CPUHard:     realTimeLimit,
		FileSize:    outputLimit << 20,
		Stack:       stackLimit << 20,
		Data:        memoryLimit << 20,
		OpenFile:    256,
		DisableCore: true,
	}
	printLimit(&rlimits)

	actionDefault := seccomp.ActionKill
	if showDetails {
		actionDefault = seccomp.ActionTrace
	}

	seccompBuilder := seccomp.Builder{
		Allow:   allow,
		Trace:   trace,
		Default: actionDefault,
	}
	filter, err := seccompBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create seccomp filter %v", err)
	}

	limit := runner.Limit{
		TimeLimit:   time.Duration(timeLimit) * time.Second,
		MemoryLimit: runner.Size(memoryLimit << 20),
	}

	r := &sandbox.Runner{
		Args:        args,
		Env:         []string{pathEnv},
		ExecFile:    execFile,
		WorkDir:     workDir,
		RLimits:     rlimits.PrepareRLimit(),
		Limit:       limit,
		Files:       fds,
		Seccomp:     filter,
		ShowDetails: showDetails,
		Unsafe:      unsafe,
		Handler:     h,
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	// Run tracer
	sTime := time.Now()
	c, cancel := context.WithTimeout(context.Background(), time.Duration(int64(realTimeLimit)*int64(time.Second)))
	defer cancel()

	s := make(chan runner.Result, 1)
	go func() {
		s <- r.Run(c)
	}()
	rTime := time.Now()

	select {
	case <-sig:
		cancel()
		rt = <-s
		rt.Status = runner.StatusSystemError

	case rt = <-s:
	}
	eTime := time.Now()

	if rt.SetUpTime == 0 {
		rt.SetUpTime = rTime.Sub(sTime)
		rt.RunningTime = eTime.Sub(rTime)
	}

	printResult(&rt)
	return &rt, nil
}
