package ptrace

import (
	"os"
	"syscall"
)

type Context struct {
	Pid  int
	regs syscall.PtraceRegs
}

var (
	UseVMReadv = true
	pageSize   = 4 << 10
)

func init() {
	pageSize = os.Getpagesize()
}

func getTrapContext(pid int) (*Context, error) {
	var regs syscall.PtraceRegs
	//err := syscall.PtraceGetRegs(pid, &regs)
	err := ptraceGetRegSet(pid, &regs)
	if err != nil {
		return nil, err
	}
	return &Context{
		Pid:  pid,
		regs: regs,
	}, nil
}


