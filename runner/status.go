package runner

type Status int

const (
	StatusNormal Status = iota // 0 normal

	StatusInvalid //1 invalid

	// Limit Exceeded
	StatusTimeLimitExceeded   // 2 TLE
	StatusMemoryLimitExceeded // 3 MLE
	StatusOutputLimitExceeded // 4 OLE

	// Syscall
	StatusDisallowedSyscall // 5 ban

	// Runtime Error
	StatusSignalled         // 6 signalled
	StatusNonzeroExitStatus // 7 nonzero exit status

	// System Error
	StatusSystemError // 8 system error
)

var (
	statusString = []string{
		"",
		"Invalid",
		"Time Limit Exceeded",
		"Memory Limit Exceeded",
		"Output Limit Exceeded",
		"Disallowed Syscall",
		"Signalled",
		"Nonzero Exit Status",
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

func (t Status) Error() string {
	return t.String()
}
