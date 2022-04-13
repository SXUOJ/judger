package runner

import (
	"fmt"
	"time"
)

type Limit struct {
	TimeLimit   time.Duration
	MemoryLimit Size
}

func (l Limit) String() string {
	return fmt.Sprintf("Limit: \nTime: %v \nMemory: %v\n", l.TimeLimit, l.MemoryLimit.String())
}
