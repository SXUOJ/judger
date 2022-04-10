package forkexec

import (
	"fmt"
	"syscall"
)

type ForkError int

type ChildError struct {
	Err     syscall.Errno
	ForkErr ForkError
	Index   int
}

const (
	FErrorClone ForkError = iota + 1
	FErrorCloseWrite
	FErrorUnshareUserRead
	FErrorGetPid
	FErrorDup3
	FErrorFcntl
	FErrorSetSid
	FErrorChdir
	FErrorSetRlimit
	FErrorSetNoNewPrivs
	FErrorPtraceMe
	FErrorStop
	FErrorSeccomp
	FErrorSyncWrite
	FErrorSyncRead
	FErrorExecve
)

var locToString = []string{
	"unknown",
	"clone",
	"close_write",
	"unshare_user_read",
	"getpid",
	"dup3",
	"fcntl",
	"setsid",
	"chdir",
	"setrlimt",
	"set_no_new_privs",
	"ptrace_me",
	"stop",
	"seccomp",
	"sync_write",
	"sync_read",
	"execve",
}

func (e ForkError) String() string {
	if e >= FErrorClone && e <= FErrorExecve {
		return locToString[e]
	}
	return "unknown"
}

func (e ChildError) Error() string {
	if e.Index > 0 {
		return fmt.Sprintf("%s(%d): %s", e.ForkErr.String(), e.Index, e.Err.Error())
	}
	return fmt.Sprintf("%s: %s", e.ForkErr.String(), e.Err.Error())
}
