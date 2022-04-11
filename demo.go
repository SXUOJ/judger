package main

import (
	"os"

	"github.com/Sxu-Online-Judge/judge/pkg/forkexec"
	"github.com/Sxu-Online-Judge/judge/pkg/rlimit"
	"github.com/Sxu-Online-Judge/judge/pkg/seccomp"
	"golang.org/x/sys/unix"
)

var (
	timeLimit      = uint64(10)   //uint64
	realTimeLimit  = uint64(10)   //uint64
	memoryLimit    = uint64(128)  //uint64
	outputLimit    = uint64(128)  //uint64
	stackLimit     = uint64(128)  //uint64
	inputFileName  = "input.txt"  //string
	outputFileName = "output.txt" //string
	errorFileName  = "error.txt"  //string

	workPath string
)

func main() {
	// execFile, err := os.Open(os.Args[1])
	// must(err)

	if workPath == "" {
		workPath, _ = os.Getwd()
	}

	files, err := prepareFiles(inputFileName, outputFileName, errorFileName)
	if err != nil {
		must(err)
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

	seccompBuilder := seccomp.Builder{
		Default: seccomp.ActionAllow,
		Trace:   []string{
			// "write",
		},
	}
	filter, err := seccompBuilder.Build()
	must(err)

	rlims := rlimit.RLimits{
		CPU:         timeLimit,
		CPUHard:     realTimeLimit,
		FileSize:    outputLimit << 20,
		Stack:       stackLimit << 20,
		Data:        memoryLimit << 20,
		OpenFile:    256,
		DisableCore: true,
	}

	ch := &forkexec.Runner{
		Args: os.Args[1:],
		Env:  os.Environ(),
		// ExecFile: execFile.Fd(),
		RLimits: rlims.PrepareRLimit(),
		Files:   fds,
		WorkDir: workPath,
		Ptrace:  true,
		Seccomp: filter.SockFprog(),
	}

	pid, err := ch.Start()
	must(err)
	var ws unix.WaitStatus
	_, err = unix.Wait4(pid, &ws, 0, nil)
	must(err)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
