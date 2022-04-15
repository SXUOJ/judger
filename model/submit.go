package model

import (
	"time"

	"github.com/SXUOJ/judge/pkg/rlimit"
	"github.com/SXUOJ/judge/runner"
	"github.com/SXUOJ/judge/worker"
	"github.com/sirupsen/logrus"
)

type Submit struct {
	SubmitID  string
	ProblemID string

	Type      string `json:"type"`
	AllowProc bool   `json:"allow_proc"`

	// Limit
	TimeLimit     uint64
	RealTimeLimit uint64
	MemoryLimit   uint64
	OutputLimit   uint64
	StackLimit    uint64
}

func (submit *Submit) Load() (*worker.Worker, error) {
	if submit.RealTimeLimit < submit.TimeLimit {
		submit.RealTimeLimit = submit.TimeLimit + 2
	}

	if submit.StackLimit > submit.MemoryLimit {
		submit.StackLimit = submit.MemoryLimit
	}

	rlimits := rlimit.RLimits{
		CPU:         submit.TimeLimit,
		CPUHard:     submit.RealTimeLimit,
		FileSize:    submit.OutputLimit << 20,
		Stack:       submit.StackLimit << 20,
		Data:        submit.MemoryLimit << 20,
		OpenFile:    256,
		DisableCore: true,
	}
	printLimit(&rlimits)

	limit := runner.Limit{
		TimeLimit:   time.Duration(submit.TimeLimit) * time.Second,
		MemoryLimit: runner.Size(submit.MemoryLimit << 20),
	}

	return &worker.Worker{
		ProblemID:     submit.ProblemID,
		SubmitID:      submit.SubmitID,
		Type:          submit.Type,
		RLimits:       rlimits.PrepareRLimit(),
		Limit:         limit,
		RealTimeLimit: submit.RealTimeLimit,
		AllowProc:     submit.AllowProc,
	}, nil
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
