package sandbox

import "syscall"

type Context struct {
	Pid  int
	regs syscall.PtraceRegs
}
