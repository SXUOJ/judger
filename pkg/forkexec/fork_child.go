package forkexec

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

func forkAndExecInChild(r *Runner, argv0 *byte, argv, env []*byte, workdir *byte, p [2]int) (r1 uintptr, err1 syscall.Errno) {
	var (
		idx        int
		pipe       int
		forkErr    ForkError
		childError ChildError
	)
	fd, nextfd := prepareFds(r.Files)

	syscall.ForkLock.Lock()
	beforeFork()
	r1, _, err1 = syscall.RawSyscall6(syscall.SYS_CLONE, uintptr(syscall.SIGCHLD)|(r.CloneFlags&UnshareFlags), 0, 0, 0, 0, 0)
	if err1 != 0 || r1 != 0 {
		fmt.Println("in parent process")
		return
	}
	afterForkInChild()

	pipe, forkErr, idx, err1 = forkAndExecInChild1(r, argv0, argv, env, workdir, fd, nextfd, p)

	childError.Err = err1
	childError.ForkErr = forkErr
	childError.Index = idx

	// send error code on pipe
	syscall.RawSyscall(unix.SYS_WRITE, uintptr(pipe), uintptr(unsafe.Pointer(&childError)), unsafe.Sizeof(childError))
	for {
		syscall.RawSyscall(syscall.SYS_EXIT, uintptr(err1), 0, 0)
	}
	return
}
