package main

import (
	"fmt"
	"os"

	"github.com/Sxu-Online-Judge/judge/pkg/rlimit"
	"github.com/Sxu-Online-Judge/judge/runner"
)

const (
	pathEnv = "PATH=/usr/local/bin:/usr/bin:/bin"
)

type Status int

const (
	StatusNormal  Status = iota // 0
	StatusInvalid               // 1
	StatusRE                    // 2
	StatusMLE                   // 3
	StatusTLE                   // 4
	StatusOLE                   // 5
	StatusBan                   // 6
	StatusFatal                 // 7
)

func getStatus(s runner.Status) int {
	switch s {
	case runner.StatusNormal:
		return int(StatusNormal)
	case runner.StatusInvalid:
		return int(StatusInvalid)
	case runner.StatusTimeLimitExceeded:
		return int(StatusTLE)
	case runner.StatusMemoryLimitExceeded:
		return int(StatusMLE)
	case runner.StatusOutputLimitExceeded:
		return int(StatusOLE)
	case runner.StatusDisallowedSyscall:
		return int(StatusBan)
	case runner.StatusSignalled, runner.StatusNonzeroExitStatus:
		return int(StatusRE)
	default:
		return int(StatusFatal)
	}
}

func printLimit(rl *rlimit.RLimits) {
	if showPrint {
		fmt.Println("cpu: ", rl.CPU)
		fmt.Println("cpuHard: ", rl.CPUHard)
		fmt.Println("fileSize: ", rl.FileSize)
		fmt.Println("stack: ", rl.Stack)
		fmt.Println("data: ", rl.Data)
		fmt.Println("openfile: ", rl.OpenFile)
		fmt.Println("disableCore", rl.DisableCore)
	}
}

func printResult(rt *runner.Result) {
	if showPrint {
		fmt.Println("status: ", rt.Status)
		fmt.Println("exitCode: ", rt.ExitCode)
		fmt.Println("error: ", rt.Error)
		fmt.Println("time: ", rt.Time)
		fmt.Println("memory: ", rt.Memory)

		fmt.Println("runTime: ", rt.RunningTime)
		fmt.Println("setUpTime: ", rt.SetUpTime)
	}
}

// func printUsage() {
// 	fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] <args>\n", os.Args[0])
// 	flag.PrintDefaults()
// 	os.Exit(2)
// }

func debug(v ...interface{}) {
	if showDetails {
		fmt.Fprintln(os.Stderr, v...)
	}
}

// prepareFile opens file for new process
func prepareFiles(inputFile, outputFile, errorFile string) ([]*os.File, error) {
	var err error
	files := make([]*os.File, 3)
	if inputFile != "" {
		files[0], err = os.OpenFile(inputFile, os.O_RDONLY, 0755)
		if err != nil {
			goto openerror
		}
	}
	if outputFile != "" {
		files[1], err = os.OpenFile(outputFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
		if err != nil {
			goto openerror
		}
	}
	if errorFile != "" {
		files[2], err = os.OpenFile(errorFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
		if err != nil {
			goto openerror
		}
	}
	return files, nil
openerror:
	closeFiles(files)
	return nil, err
}

// closeFiles close all file in the list
func closeFiles(files []*os.File) {
	for _, f := range files {
		if f != nil {
			f.Close()
		}
	}
}
