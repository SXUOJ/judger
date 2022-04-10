package forkexec

import "syscall"

func forkAndExecInChild(r *Runner, argv0 *byte, argv, env []*byte, workdir *byte, p [2]int) (r1 uintptr, err1 syscall.Errno) {

	return
}
