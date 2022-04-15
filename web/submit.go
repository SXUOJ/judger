package web

import (
	"github.com/SXUOJ/judge/worker"
)

type Submit struct {
	SubmitID  string `json:"submit_id"`
	ProblemID string `json:"problem_id"`

	FileName string `json:"file_name"`
	Type     string `json:"type"`

	AllowProc bool `json:"allow_proc"`

	// Limit
	TimeLimit     uint64 `json:"time_limit"`
	RealTimeLimit uint64 `json:"real_time_limit"`
	MemoryLimit   uint64 `json:"memory_limit"`
	OutputLimit   uint64 `json:"output_limit"`
	StackLimit    uint64 `json:"stack_limit"`
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

		FileName:  submit.FileName,
		Type:      submit.Type,
		AllowProc: submit.AllowProc,

		TimeLimit:     submit.TimeLimit,
		RealTimeLimit: submit.RealTimeLimit,
		MemoryLimit:   submit.MemoryLimit,
		OutputLimit:   submit.OutputLimit,
		StackLimit:    submit.StackLimit,
	}, nil
}
