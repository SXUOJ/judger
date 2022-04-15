package worker

import (
	"time"

	"github.com/SXUOJ/judge/runner"
)

type Results interface{}

type CompileResult Result

type RunResults []RunResult

type RunResult struct {
	SampleId int `json:"sample_id"`
	Result
}

type Result struct {
	Status Status `json:"status"`

	SetUpTime   time.Duration `json:"set_up_time"`
	RunningTime time.Duration `json:"running_time"`
	Time        time.Duration `json:"time"`
	Memory      runner.Size   `json:"memory"`

	ExitCode int    `json:"exit_code"`
	Error    string `json:"error"`
}
