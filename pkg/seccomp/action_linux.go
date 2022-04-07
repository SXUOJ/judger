package seccomp

import (
	seccomp "github.com/elastic/go-seccomp-bpf"
)

func ToSeccompAction(a Action) seccomp.Action {
	var action seccomp.Action
	switch a.Action() {
	case ActionAllow:
		action = seccomp.ActionAllow
	case ActionErrno:
		action = seccomp.ActionErrno
	case ActionTrace:
		action = seccomp.ActionTrace
	default:
		action = seccomp.ActionKillProcess
	}
	return action
}
