package handle

import (
	"github.com/Sxu-Online-Judge/judge/ptrace"
)

// Handler defines file access restricted handler to call the ptrace
// safe runner
type Handler struct {
	SyscallCounter SyscallCounter
}

// CheckSyscall checks syscalls other than allowed and traced against the
// SyscallCounter
func (h *Handler) CheckSyscall(syscallName string) ptrace.TraceAction {
	// if it is traced, then try to count syscall
	if inside, allow := h.SyscallCounter.Check(syscallName); inside {
		if allow {
			return ptrace.ActionAllow
		}
		return ptrace.ActionKill
	}
	// if it is traced but not counted, it should be soft banned
	return ptrace.ActionTrace
}
