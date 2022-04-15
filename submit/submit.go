package submit

type Submit struct {
	SubmitID  string
	ProblemID string

	Type      string `json:"type"`
	AllowProc bool   `json:"allow_proc"`

	WorkDir string

	// Limit
	TimeLimit     uint64
	RealTimeLimit uint64
	MemoryLimit   uint64
	OutPutLimit   uint64
	StackLimit    uint64
}


