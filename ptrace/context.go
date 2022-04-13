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

func (c *Context) GetString(addr uintptr) string {
	buff := make([]byte, syscall.PathMax)
	if UseVMReadv {
		if err := vmReadStr(c.Pid, addr, buff); err != nil {
			// if ENOSYS, then disable this function
			if no, ok := err.(syscall.Errno); ok {
				if no == syscall.ENOSYS {
					UseVMReadv = false
				}
			}
		} else {
			return string(buff[:clen(buff)])
		}
	}
	syscall.PtracePeekData(c.Pid, addr, buff)
	return string(buff[:clen(buff)])
}
