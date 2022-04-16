package worker

import "github.com/SXUOJ/judge/runner"

type Status int

const (
	StatusNormal Status = iota
	StatusAC
	StatusWA
	StatusCE

	StatusRE //runtime error

	// Limit Exceeded
	StatusTLE
	StatusMLE
	StatusOLE

	//presentation error
	StatusPE

	StatusSE //system error
)

var (
	statusString = []string{
		"",
		"Accepted",
		"Wrong Answer",
		"Compile Error",
		"Runtime Error",
		"Time Limit Exceed",
		"Memory Limit Exceed",
		"Output Limit Exceed",
		"Presentation Error",
		"System Error",
	}
)

func (t Status) String() string {
	i := int(t)
	if i >= 0 && i < len(statusString) {
		return statusString[i]
	}
	return statusString[0]
}

func convertStatus(status runner.Status) Status {
	switch status {
	case runner.StatusNormal:
		return StatusNormal
	case runner.StatusInvalid,
		runner.StatusDisallowedSyscall,
		runner.StatusSignalled,
		runner.StatusNonzeroExitStatus:
		return StatusRE
	case runner.StatusTimeLimitExceeded:
		return StatusTLE
	case runner.StatusMemoryLimitExceeded:
		return StatusMLE
	case runner.StatusOutputLimitExceeded:
		return StatusOLE
	case runner.StatusSystemError:
		return StatusSE
	default:
		return StatusSE
	}
}
