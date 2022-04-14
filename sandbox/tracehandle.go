package sandbox

import (
	"fmt"
	"os"
	"syscall"

	"github.com/SXUOJ/judge/pkg/seccomp"
	"github.com/SXUOJ/judge/ptrace"
)

var (
	BanRet = syscall.EACCES
)

type tracerHandler struct {
	ShowDetails, Unsafe bool
	Handler             Handler
}

func (h *tracerHandler) Debug(v ...interface{}) {
	if h.ShowDetails {
		fmt.Fprintln(os.Stderr, v...)
	}
}

func (h *tracerHandler) Handle(ctx *ptrace.Context) ptrace.TraceAction {
	syscallNo := ctx.SyscallNo()
	syscallName, err := seccomp.ToSyscallName(syscallNo)
	h.Debug("syscall:", syscallNo, syscallName, err)
	if err != nil {
		h.Debug("invalid syscall no")
		return ptrace.ActionKill
	}

	action := ptrace.ActionKill
	action = h.Handler.CheckSyscall(syscallName)
	if h.Unsafe && action == ptrace.ActionKill {
		action = ptrace.ActionTrace
	}

	switch action {
	case ptrace.ActionAllow:
		return ptrace.ActionAllow
	case ptrace.ActionTrace:
		h.Debug("<soft ban syscall>")
		return softTraceSyscall(ctx)
	default:
		return ptrace.ActionKill
	}
}

func softTraceSyscall(ctx *ptrace.Context) ptrace.TraceAction {
	ctx.SetReturnValue(-int(BanRet))
	return ptrace.ActionTrace
}
