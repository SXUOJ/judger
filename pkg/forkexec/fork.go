package forkexec

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
	// required for go:linkname.
)

func (r *Runner) Start() (int, error) {
	argv0, argv, env, err := prepareExec(r.Args, r.Env)
	if err != nil {
		return 0, err
	}

	// prepare work dir
	workdir, err := syscallStringFromString(r.WorkDir)
	if err != nil {
		return 0, err
	}

	p, err := syscall.Socketpair(syscall.AF_LOCAL, syscall.SOCK_STREAM|syscall.SOCK_CLOEXEC, 0)
	if err != nil {
		return 0, err
	}

	pid, err1 := forkAndExecInChild(r, argv0, argv, env, workdir, p)

	afterFork()
	syscall.ForkLock.Unlock()
	return syncWithChild(r, p, int(pid), err1)
}

func syncWithChild(r *Runner, p [2]int, pid int, err1 syscall.Errno) (int, error) {
	var (
		r1       uintptr
		err2     syscall.Errno
		err      error
		childErr ChildError
	)

	// sync with child
	unix.Close(p[1])

	// clone syscall failed
	if err1 != 0 {
		unix.Close(p[0])
		childErr.ForkErr = FErrorClone
		childErr.Err = err1
		return 0, childErr
	}

	r1, _, err1 = syscall.RawSyscall(syscall.SYS_READ, uintptr(p[0]), uintptr(unsafe.Pointer(&childErr)), uintptr(unsafe.Sizeof(childErr)))
	// child returned error code
	if (r1 != unsafe.Sizeof(err2) && r1 != unsafe.Sizeof(childErr)) || childErr.Err != 0 || err1 != 0 {
		childErr.Err = handlePipeError(r1, childErr.Err)
		goto fail
	}

	// if syncfunc return error, then fail child immediately
	if r.SyncFunc != nil {
		if err = r.SyncFunc(int(pid)); err != nil {
			goto fail
		}
	}
	// otherwise, ack child (err1 == 0)
	syscall.RawSyscall(syscall.SYS_WRITE, uintptr(p[0]), uintptr(unsafe.Pointer(&err1)), uintptr(unsafe.Sizeof(err1)))

	// if stopped before execve by signal SIGSTOP or PTRACE_ME, then do not wait until execve
	if r.Ptrace {
		// let's wait it in another goroutine to avoid SIGPIPE
		go func() {
			syscall.RawSyscall(syscall.SYS_READ, uintptr(p[0]), uintptr(unsafe.Pointer(&childErr)), uintptr(unsafe.Sizeof(childErr)))
			unix.Close(p[0])
		}()
		return int(pid), nil
	}

	// if read anything mean child failed after sync (close_on_exec so it should not block)
	r1, _, err1 = syscall.RawSyscall(syscall.SYS_READ, uintptr(p[0]), uintptr(unsafe.Pointer(&childErr)), uintptr(unsafe.Sizeof(childErr)))
	unix.Close(p[0])
	if r1 != 0 || err1 != 0 {
		childErr.Err = handlePipeError(r1, childErr.Err)
		goto failAfterClose
	}
	return int(pid), nil

fail:
	unix.Close(p[0])

failAfterClose:
	handleChildFailed(int(pid))
	if childErr.Err == 0 {
		return 0, err
	}
	return 0, childErr
}

// check pipe error
func handlePipeError(r1 uintptr, errno syscall.Errno) syscall.Errno {
	if r1 >= unsafe.Sizeof(errno) {
		return syscall.Errno(errno)
	}
	return syscall.EPIPE
}

func handleChildFailed(pid int) {
	var wstatus syscall.WaitStatus
	// make sure not blocked
	syscall.Kill(pid, syscall.SIGKILL)
	// child failed; wait for it to exit, to make sure the zombies don't accumulate
	_, err := syscall.Wait4(pid, &wstatus, 0, nil)
	for err == syscall.EINTR {
		_, err = syscall.Wait4(pid, &wstatus, 0, nil)
	}
}
