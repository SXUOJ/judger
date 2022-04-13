package runner

import (
	"time"
)

type Result struct {
	Status
	ExitCode int
	Error    string

	SetUpTime   time.Duration
	RunningTime time.Duration
	Time        time.Duration
	Memory      Size
}
