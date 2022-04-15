package worker

type Status int

const (
	_ Status = iota
	StatusAC
	StatusWA
	StatusCE
	StatusRE
	StatusTLE
	StatusMLE
	StatusOLE
	StatusPE
	StatusSE
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
