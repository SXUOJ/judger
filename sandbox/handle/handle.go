package handle

import (
	"github.com/Sxu-Online-Judge/judge/ptrace"
)

// Handler defines file access restricted handler to call the ptrace
// safe runner
type Handler struct {
	FileSet        *FileSets
	SyscallCounter SyscallCounter
}

// CheckRead checks whether the file have read permission
func (h *Handler) CheckRead(fn string) ptrace.TraceAction {
	if !h.FileSet.IsReadableFile(fn) {
		return h.onDgsFileDetect(fn)
	}
	return ptrace.ActionAllow
}

// CheckWrite checks whether the file have write permission
func (h *Handler) CheckWrite(fn string) ptrace.TraceAction {
	if !h.FileSet.IsWritableFile(fn) {
		return h.onDgsFileDetect(fn)
	}
	return ptrace.ActionAllow
}

// CheckStat checks whether the file have stat permission
func (h *Handler) CheckStat(fn string) ptrace.TraceAction {
	if !h.FileSet.IsStatableFile(fn) {
		return h.onDgsFileDetect(fn)
	}
	return ptrace.ActionAllow
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

// onDgsFileDetect soft ban file if in soft ban set
// otherwise stops the trace process
func (h *Handler) onDgsFileDetect(name string) ptrace.TraceAction {
	if h.FileSet.IsSoftBanFile(name) {
		return ptrace.ActionTrace
	}
	return ptrace.ActionKill
}
