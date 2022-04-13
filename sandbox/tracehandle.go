package sandbox

import (
	"fmt"
	"os"
	"syscall"

	"github.com/Sxu-Online-Judge/judge/pkg/seccomp"
	"github.com/Sxu-Online-Judge/judge/ptrace"
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

// func (h *tracerHandler) getString(ctx *ptrace.Context, addr uint) string {
// 	return absPath(ctx.Pid, ctx.GetString(uintptr(addr)))
// }

// func (h *tracerHandler) checkOpen(ctx *ptrace.Context, addr uint, flags uint) ptrace.TraceAction {
// 	fn := h.getString(ctx, addr)
// 	isReadOnly := (flags&syscall.O_ACCMODE == syscall.O_RDONLY) &&
// 		(flags&syscall.O_CREAT == 0) &&
// 		(flags&syscall.O_EXCL == 0) &&
// 		(flags&syscall.O_TRUNC == 0)

// 	h.Debug("open: ", fn, getFileMode(flags))
// 	if isReadOnly {
// 		return h.Handler.CheckRead(fn)
// 	}
// 	return h.Handler.CheckWrite(fn)
// }

// func (h *tracerHandler) checkRead(ctx *ptrace.Context, addr uint) ptrace.TraceAction {
// 	fn := h.getString(ctx, addr)
// 	h.Debug("check read: ", fn)
// 	return h.Handler.CheckRead(fn)
// }

// func (h *tracerHandler) checkWrite(ctx *ptrace.Context, addr uint) ptrace.TraceAction {
// 	fn := h.getString(ctx, addr)
// 	h.Debug("check write: ", fn)
// 	return h.Handler.CheckWrite(fn)
// }

// func (h *tracerHandler) checkStat(ctx *ptrace.Context, addr uint) ptrace.TraceAction {
// 	fn := h.getString(ctx, addr)
// 	h.Debug("check stat: ", fn)
// 	return h.Handler.CheckStat(fn)
// }

func (h *tracerHandler) Handle(ctx *ptrace.Context) ptrace.TraceAction {
	syscallNo := ctx.SyscallNo()
	syscallName, err := seccomp.ToSyscallName(syscallNo)
	h.Debug("syscall:", syscallNo, syscallName, err)
	if err != nil {
		h.Debug("invalid syscall no")
		return ptrace.ActionKill
	}

	action := ptrace.ActionKill
	// switch syscallName {
	// case "open":
	// 	action = h.checkOpen(ctx, ctx.Arg0(), ctx.Arg1())
	// case "openat":
	// 	action = h.checkOpen(ctx, ctx.Arg1(), ctx.Arg2())

	// case "readlink":
	// 	action = h.checkRead(ctx, ctx.Arg0())
	// case "readlinkat":
	// 	action = h.checkRead(ctx, ctx.Arg1())

	// case "unlink":
	// 	action = h.checkWrite(ctx, ctx.Arg0())
	// case "unlinkat":
	// 	action = h.checkWrite(ctx, ctx.Arg1())

	// case "access":
	// 	action = h.checkStat(ctx, ctx.Arg0())
	// case "faccessat", "newfstatat":
	// 	action = h.checkStat(ctx, ctx.Arg1())

	// case "stat", "stat64":
	// 	action = h.checkStat(ctx, ctx.Arg0())
	// case "lstat", "lstat64":
	// 	action = h.checkStat(ctx, ctx.Arg0())

	// case "execve":
	// 	action = h.checkRead(ctx, ctx.Arg0())
	// case "execveat":
	// 	action = h.checkRead(ctx, ctx.Arg1())

	// case "chmod":
	// 	action = h.checkWrite(ctx, ctx.Arg0())
	// case "rename":
	// 	action = h.checkWrite(ctx, ctx.Arg0())

	// default:
	action = h.Handler.CheckSyscall(syscallName)
	if h.Unsafe && action == ptrace.ActionKill {
		action = ptrace.ActionTrace
	}
	// }

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

// func getFileMode(flags uint) string {
// 	switch flags & syscall.O_ACCMODE {
// 	case syscall.O_RDONLY:
// 		return "r "
// 	case syscall.O_WRONLY:
// 		return "w "
// 	case syscall.O_RDWR:
// 		return "wr"
// 	default:
// 		return "??"
// 	}
// }

// getProcCwd gets the process CWD
func getProcCwd(pid int) string {
	fileName := "/proc/self/cwd"
	if pid > 0 {
		fileName = fmt.Sprintf("/proc/%d/cwd", pid)
	}
	s, err := os.Readlink(fileName)
	if err != nil {
		return ""
	}
	return s
}

// absPath calculates the absolute path for a process
// built-in function did the dirty works to resolve relative paths
// func absPath(pid int, p string) string {
// 	// if relative path
// 	if !path.IsAbs(p) {
// 		return path.Join(getProcCwd(pid), p)
// 	}
// 	return path.Clean(p)
// }
