package runner

import (
	"time"
)

type Result struct {
	Status
	Signal int
	Error  string

	SetUpTime   time.Duration
	RunningTime time.Duration
	Time        time.Duration
	Memory      Size
}
