package forkexec

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

// Reference to src/syscall/exec_linux.go
//go:norace
func forkAndExecInChild(r *Runner, argv0 *byte, argv, env []*byte, workdir *byte, p [2]int) (r1 uintptr, err1 syscall.Errno) {
	var (
		loc        ForkError
		idx, pipe  int
		childError ChildError
	)
	fd, nextfd := prepareFds(r.Files)

	syscall.ForkLock.Lock()
	beforeFork()
	r1, _, err1 = syscall.RawSyscall6(syscall.SYS_CLONE, uintptr(syscall.SIGCHLD), 0, 0, 0, 0, 0)
	if err1 != 0 || r1 != 0 {
		return
	}
	afterForkInChild()

	pipe, loc, idx, err1 = forkAndExecInChild1(r, argv0, argv, env, workdir, fd, nextfd, p)

	// send error code on pipe
	childError.Err = err1
	childError.ForkErr = loc
	childError.Index = idx

	// send error code on pipe
	syscall.RawSyscall(unix.SYS_WRITE, uintptr(pipe), uintptr(unsafe.Pointer(&childError)), unsafe.Sizeof(childError))
	for {
		syscall.RawSyscall(syscall.SYS_EXIT, uintptr(err1), 0, 0)
	}
}

//go:nosplit
func forkAndExecInChild1(r *Runner, argv0 *byte, argv, env []*byte, workdir *byte, fd []int, nextfd int, p [2]int) (pipe int, loc ForkError, idx int, err1 syscall.Errno) {
	pipe = p[1]
	var (
		r1, pid     uintptr
		err2        syscall.Errno
		unshareUser = r.CloneFlags&unix.CLONE_NEWUSER == unix.CLONE_NEWUSER
	)

	// Close write end of pipe
	if _, _, err1 = syscall.RawSyscall(syscall.SYS_CLOSE, uintptr(p[0]), 0, 0); err1 != 0 {
		return pipe, FErrorCloseWrite, 0, err1
	}

	if unshareUser {
		r1, _, err1 = syscall.RawSyscall(syscall.SYS_READ, uintptr(pipe), uintptr(unsafe.Pointer(&err2)), unsafe.Sizeof(err2))
		if err1 != 0 {
			return pipe, FErrorUnshareUserRead, 0, err1
		}
		if r1 != unsafe.Sizeof(err2) {
			err1 = syscall.EINVAL
			return pipe, FErrorUnshareUserRead, 0, err1
		}
		if err2 != 0 {
			err1 = err2
			return pipe, FErrorUnshareUserRead, 0, err1
		}
	}

	// Get pid of child
	pid, _, err1 = syscall.RawSyscall(syscall.SYS_GETPID, 0, 0, 0)
	if err1 != 0 {
		return pipe, FErrorGetPid, 0, err1
	}

	if pipe < nextfd {
		_, _, err1 = syscall.RawSyscall(syscall.SYS_DUP3, uintptr(pipe), uintptr(nextfd), syscall.O_CLOEXEC)
		if err1 != 0 {
			return pipe, FErrorDup3, 0, err1
		}
		pipe = nextfd
		nextfd++
	}
	if r.ExecFile > 0 && int(r.ExecFile) < nextfd {
		_, _, err1 = syscall.RawSyscall(syscall.SYS_DUP3, r.ExecFile, uintptr(nextfd), syscall.O_CLOEXEC)
		if err1 != 0 {
			return pipe, FErrorDup3, 0, err1
		}
		r.ExecFile = uintptr(nextfd)
		nextfd++
	}
	for i := 0; i < len(fd); i++ {
		if fd[i] >= 0 && fd[i] < int(i) {
			// Avoid fd rewrite
			for nextfd == pipe || (r.ExecFile > 0 && nextfd == int(r.ExecFile)) {
				nextfd++
			}
			_, _, err1 = syscall.RawSyscall(syscall.SYS_DUP3, uintptr(fd[i]), uintptr(nextfd), syscall.O_CLOEXEC)
			if err1 != 0 {
				return pipe, FErrorDup3, 0, err1
			}
			// Set up close on exec
			fd[i] = nextfd
			nextfd++
		}
	}
	for i := 0; i < len(fd); i++ {
		if fd[i] == -1 {
			syscall.RawSyscall(syscall.SYS_CLOSE, uintptr(i), 0, 0)
			continue
		}
		if fd[i] == int(i) {
			// dup2(i, i) will not clear close on exec flag, need to reset the flag
			_, _, err1 = syscall.RawSyscall(syscall.SYS_FCNTL, uintptr(fd[i]), syscall.F_SETFD, 0)
			if err1 != 0 {
				return pipe, FErrorFcntl, 0, err1
			}
			continue
		}
		_, _, err1 = syscall.RawSyscall(syscall.SYS_DUP3, uintptr(fd[i]), uintptr(i), 0)
		if err1 != 0 {
			return pipe, FErrorDup3, 0, err1
		}
	}

	// set session id
	_, _, err1 = syscall.RawSyscall(syscall.SYS_SETSID, 0, 0, 0)
	if err1 != 0 {
		return pipe, FErrorSetSid, 0, err1
	}

	//  chdir
	if workdir != nil {
		_, _, err1 = syscall.RawSyscall(syscall.SYS_CHDIR, uintptr(unsafe.Pointer(workdir)), 0, 0)
		if err1 != 0 {
			return pipe, FErrorChdir, 0, err1
		}
	}

	// set rlimit
	for i, rlim := range r.RLimits {
		_, _, err1 = syscall.RawSyscall6(syscall.SYS_PRLIMIT64, 0, uintptr(rlim.Res), uintptr(unsafe.Pointer(&rlim.Rlim)), 0, 0, 0)
		if err1 != 0 {
			return pipe, FErrorSetRlimit, i, err1
		}
	}

	// No new privs
	if r.Seccomp != nil {
		_, _, err1 = syscall.RawSyscall6(syscall.SYS_PRCTL, unix.PR_SET_NO_NEW_PRIVS, 1, 0, 0, 0, 0)
		if err1 != 0 {
			return pipe, FErrorSetNoNewPrivs, 0, err1
		}
	}

	if r.Ptrace && r.Seccomp != nil {
		{
			r1, _, err1 = syscall.RawSyscall(syscall.SYS_WRITE, uintptr(pipe), uintptr(unsafe.Pointer(&err2)), uintptr(unsafe.Sizeof(err2)))
			if r1 == 0 || err1 != 0 {
				return pipe, FErrorSyncWrite, 0, err1
			}

			r1, _, err1 = syscall.RawSyscall(syscall.SYS_READ, uintptr(pipe), uintptr(unsafe.Pointer(&err2)), uintptr(unsafe.Sizeof(err2)))
			if r1 == 0 || err1 != 0 {
				return pipe, FErrorSyncRead, 0, err1
			}
		}
		_, _, err1 = syscall.RawSyscall(syscall.SYS_PTRACE, uintptr(syscall.PTRACE_TRACEME), 0, 0)
		if err1 != 0 {
			return pipe, FErrorPtraceMe, 0, err1
		}
	}

	if r.Seccomp != nil && r.Ptrace {
		// Stop to wait for ptrace tracer
		_, _, err1 = syscall.RawSyscall(syscall.SYS_KILL, pid, uintptr(syscall.SIGSTOP), 0)
		if err1 != 0 {
			return pipe, FErrorStop, 0, err1
		}
	}

	// load seccomp filter
	if r.Seccomp != nil && r.Ptrace {
		_, _, err1 = syscall.RawSyscall(unix.SYS_SECCOMP, SECCOMP_SET_MODE_FILTER, SECCOMP_FILTER_FLAG_TSYNC, uintptr(unsafe.Pointer(r.Seccomp)))
		if err1 != 0 {
			return pipe, FErrorSeccomp, 0, err1
		}
	}

	if !r.Ptrace || r.Seccomp == nil {
		{
			r1, _, err1 = syscall.RawSyscall(syscall.SYS_WRITE, uintptr(pipe), uintptr(unsafe.Pointer(&err2)), uintptr(unsafe.Sizeof(err2)))
			if r1 == 0 || err1 != 0 {
				return pipe, FErrorSyncWrite, 0, err1
			}

			r1, _, err1 = syscall.RawSyscall(syscall.SYS_READ, uintptr(pipe), uintptr(unsafe.Pointer(&err2)), uintptr(unsafe.Sizeof(err2)))
			if r1 == 0 || err1 != 0 {
				return pipe, FErrorSyncRead, 0, err1
			}
		}
	}

	// Enable ptrace if no seccomp is needed
	if r.Ptrace && r.Seccomp == nil {
		_, _, err1 = syscall.RawSyscall(syscall.SYS_PTRACE, uintptr(syscall.PTRACE_TRACEME), 0, 0)
		if err1 != 0 {
			return pipe, FErrorPtraceMe, 0, err1
		}
	}

	if r.ExecFile > 0 {
		_, _, err1 = syscall.RawSyscall6(unix.SYS_EXECVEAT, r.ExecFile,
			uintptr(unsafe.Pointer(&empty[0])), uintptr(unsafe.Pointer(&argv[0])),
			uintptr(unsafe.Pointer(&env[0])), unix.AT_EMPTY_PATH, 0)
	} else {
		_, _, err1 = syscall.RawSyscall(unix.SYS_EXECVE, uintptr(unsafe.Pointer(argv0)),
			uintptr(unsafe.Pointer(&argv[0])), uintptr(unsafe.Pointer(&env[0])))
	}

	for range [50]struct{}{} {
		if err1 != syscall.ETXTBSY {
			break
		}
		syscall.RawSyscall(unix.SYS_NANOSLEEP, uintptr(unsafe.Pointer(&etxtbsyRetryInterval)), 0, 0)
		if r.ExecFile > 0 {
			_, _, err1 = syscall.RawSyscall6(unix.SYS_EXECVEAT, r.ExecFile,
				uintptr(unsafe.Pointer(&empty[0])), uintptr(unsafe.Pointer(&argv[0])),
				uintptr(unsafe.Pointer(&env[0])), unix.AT_EMPTY_PATH, 0)
		} else {
			_, _, err1 = syscall.RawSyscall(unix.SYS_EXECVE, uintptr(unsafe.Pointer(argv0)),
				uintptr(unsafe.Pointer(&argv[0])), uintptr(unsafe.Pointer(&env[0])))
		}
	}
	return pipe, FErrorExecve, 0, err1
}
