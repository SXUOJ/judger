package sandbox

type TraceAction int

const (
	ActionAllow TraceAction = iota
	ActionTrace
	ActionKill
)
