package runner

import (
	"time"
)

type Result struct {
	Status   `json:"status"`
	ExitCode int    `json:"exit_code"`
	Error    string `json:"err"`

	SetUpTime   time.Duration `json:"set_up_time"`
	RunningTime time.Duration `json:"running_time"`
	Time        time.Duration `json:"time"`
	Memory      Size          `json:"memory"`
}
