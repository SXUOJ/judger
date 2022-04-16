package worker

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
