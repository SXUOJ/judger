package forkexec

import (
	"syscall"
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
	return 0, nil
}
