package libseccomp

type Action uint32

const (
	ActionAllow Action = iota + 1
	ActionErrno
	ActionTrace
	ActionKill
)

func (a Action) Action() Action {
	return Action(a & 0xffff)
}
