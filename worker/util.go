package worker

import (
	"time"

	"github.com/SXUOJ/judge/pkg/rlimit"
	"github.com/SXUOJ/judge/runner"
	"github.com/sirupsen/logrus"
)

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

func printLimit(rl *rlimit.RLimits) {
	logrus.Debug(
		"\ncpu: ", rl.CPU,
		"\ncpuHard: ", rl.CPUHard,
		"\nfileSize: ", rl.FileSize,
		"\nstack: ", rl.Stack,
		"\ndata: ", rl.Data,
		"\nopenfile: ", rl.OpenFile,
		"\ndisableCore", rl.DisableCore,
	)
}
