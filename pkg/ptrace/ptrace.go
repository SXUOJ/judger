package ptrace

import (
	"fmt"
	"syscall"

	"github.com/Sxu-Online-Judge/judge/model"
)

var (
	AllowSysCall = []uint64{
		syscall.SYS_READ,
		syscall.SYS_WRITE,
		syscall.SYS_OPEN,
		syscall.SYS_MMAP,
		syscall.SYS_MUNMAP,
		syscall.SYS_BRK,
		syscall.SYS_WRITEV,
		syscall.SYS_ACCESS,
		syscall.SYS_EXECVE,
		syscall.SYS_UNAME,
		syscall.SYS_READLINK,
		syscall.SYS_ARCH_PRCTL,
		syscall.SYS_EXIT_GROUP,
		syscall.SYS_MQ_OPEN,
		syscall.SYS_IOPRIO_GET,
		syscall.SYS_TIME,
		syscall.SYS_READ,
		syscall.SYS_UNAME,
		syscall.SYS_WRITE,
		syscall.SYS_OPEN,
		syscall.SYS_CLOSE,
		syscall.SYS_EXECVE,
		syscall.SYS_ACCESS,
		syscall.SYS_BRK,
		syscall.SYS_MUNMAP,
		syscall.SYS_MPROTECT,
		syscall.SYS_MMAP,
		syscall.SYS_FSTAT,
		syscall.SYS_SET_THREAD_AREA,
		syscall.SYS_ARCH_PRCTL,
	}

	DangerSyscall = []uint64{
		syscall.SYS_GETPID,
	}
)

type SyscallCounter []int

const maxSyscalls = 303

func IsAllowSysCall(id uint64) model.JudgeStatus {
	if id > syscall.SYS_PRLIMIT64 {
		return model.StatusRE
	}

	for _, v := range AllowSysCall {
		if v == id {
			return model.Normal
		}
	}
	return model.Danger
}

func (s SyscallCounter) Init() SyscallCounter {
	return make(SyscallCounter, maxSyscalls)
}

func (s SyscallCounter) Inc(syscallId uint64) error {
	if syscallId > maxSyscalls {
		return fmt.Errorf("invalid syscall Id: %x", syscallId)
	}

	s[syscallId]++
	return nil
}
