package forkexec

import (
	"syscall"

	"github.com/Sxu-Online-Judge/judge/pkg/rlimit"
)

type Runner struct {
	Args     []string
	Env      []string
	ExecFile uintptr
	RLimits  []rlimit.RLimit
	Files    []uintptr
	WorkDir  string
	Ptrace   bool
	Seccomp  *syscall.SockFprog

	CloneFlags uintptr
	SyncFunc   func(int) error
}
