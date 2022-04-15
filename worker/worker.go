package worker

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/SXUOJ/judge/config"
	"github.com/SXUOJ/judge/lang"
	"github.com/SXUOJ/judge/pkg/rlimit"
	"github.com/SXUOJ/judge/pkg/seccomp"
	"github.com/SXUOJ/judge/runner"
	"github.com/SXUOJ/judge/sandbox"
)

var (
	showDetails = true
	unSafe      = true
	runDir      = filepath.Join("./tmp")
	pathEnv     = "PATH=/usr/local/bin:/usr/bin:/bin"
)

type Worker struct {
	ProblemID string
	SubmitID  string

	FileName string
	Type     string
	WorkDir  string

	AllowProc     bool
	TimeLimit     uint64
	RealTimeLimit uint64
	MemoryLimit   uint64
	OutputLimit   uint64
	StackLimit    uint64
}

func (worker *Worker) Run() (Results, error) {
	worker.WorkDir = filepath.Join(runDir, worker.SubmitID)
	var (
		result     *Result
		sourcePath = filepath.Join(worker.WorkDir, worker.FileName)
		binaryPath = filepath.Join(worker.WorkDir, worker.SubmitID)
	)

	lang, err := lang.NewLang(worker.Type, sourcePath, binaryPath)
	if err != nil {
		return nil, err
	}

	if lang.NeedCompile() {
		r, err := worker.load(-1, lang)
		if err != nil {
			return &CompileResult{
				Status: StatusCE,
				Error:  err.Error(),
			}, nil
		}
		res, err := runByOne(r, lang.CompileRealTimeLimit())
		if err != nil {
			return &CompileResult{
				Status: StatusCE,
				Error:  err.Error(),
			}, nil
		}
		if res.Status != runner.StatusNormal {
			return &CompileResult{
				Status: StatusCE,
			}, nil
		}
	}
	return result, nil
}

func (worker *Worker) load(sampleID int, lang lang.Lang) (*sandbox.Runner, error) {
	if sampleID < 0 {
		rlimits, limit := parseLimit(lang.CompileCpuTimeLimit(), lang.CompileRealTimeLimit(), 0, 0, lang.CompileMemoryLimit())

		defaultAction, allow, trace, h := config.GetConf(strings.Join([]string{worker.Type, "compile"}, "-"), worker.AllowProc)
		seccompBuilder := seccomp.Builder{
			Default: defaultAction,
			Allow:   allow,
			Trace:   trace,
		}
		filter, _ := seccompBuilder.Build()

		return &sandbox.Runner{
			Args: lang.CompileArgs(),
			Env:  []string{pathEnv},
			// ExecFile:    execFile,
			// Files:       fds,
			// WorkDir:     worker.WorkDir,
			Seccomp:     filter,
			RLimits:     rlimits.PrepareRLimit(),
			Limit:       limit,
			ShowDetails: showDetails,
			Unsafe:      unSafe,
			Handler:     h,
		}, nil
	}

	var (
		inputFileName  = ""
		outputFileName = ""
		errorFileName  = ""
	)

	defaultAction, allow, trace, h := config.GetConf(strings.Join([]string{worker.Type, "run"}, ""), worker.AllowProc)

	// limit
	rlimits, limit := parseLimit(worker.TimeLimit, worker.RealTimeLimit, worker.OutputLimit, worker.StackLimit, worker.MemoryLimit)

	// build seccomp filter
	seccompBuilder := seccomp.Builder{
		Default: defaultAction,
		Allow:   allow,
		Trace:   trace,
	}
	filter, _ := seccompBuilder.Build()

	// load fds
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

	return &sandbox.Runner{
		Args: lang.CompileArgs(),
		Env:  os.Environ(),
		// ExecFile: execFile,
		Files:       fds,
		WorkDir:     worker.WorkDir,
		Seccomp:     filter,
		RLimits:     rlimits.PrepareRLimit(),
		Limit:       limit,
		ShowDetails: showDetails,
		Unsafe:      unSafe,
		Handler:     h,
	}, nil
}

func runByOne(r *sandbox.Runner, realTimeLimit uint64) (*runner.Result, error) {
	var (
		rt runner.Result
	)

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

	printResult(rt)
	return &rt, nil
}

func parseLimit(timeLimit, realTimeLimit, outputLimit, stackLimit, memoryLimit uint64) (rlimit.RLimits, runner.Limit) {
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

	limit := runner.Limit{
		TimeLimit:   time.Duration(timeLimit) * time.Second,
		MemoryLimit: runner.Size(memoryLimit << 20),
	}

	return rlimits, limit
}
