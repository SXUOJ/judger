package forkexec

import (
	"golang.org/x/sys/unix"
)

const (
	SECCOMP_SET_MODE_STRICT   = 0
	SECCOMP_SET_MODE_FILTER   = 1
	SECCOMP_FILTER_FLAG_TSYNC = 1

	UnshareFlags = unix.CLONE_NEWIPC | unix.CLONE_NEWNET | unix.CLONE_NEWNS |
		unix.CLONE_NEWPID | unix.CLONE_NEWUSER | unix.CLONE_NEWUTS | unix.CLONE_NEWCGROUP
)

var (
	empty = []byte("\000")

	etxtbsyRetryInterval = unix.Timespec{
		Nsec: 1 * 1000 * 1000,
	}
)
