package model

import (
	"github.com/SXUOJ/judge/worker"
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

	return &worker.Worker{
		ProblemID: submit.ProblemID,
		SubmitID:  submit.SubmitID,
		Type:      submit.Type,

		RealTimeLimit: submit.RealTimeLimit,
		AllowProc:     submit.AllowProc,
	}, nil
}
