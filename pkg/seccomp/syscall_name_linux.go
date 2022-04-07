package seccomp

import (
	"fmt"

	"github.com/elastic/go-seccomp-bpf/arch"
)

var info, errInfo = arch.GetInfo("")

// ToSyscallName convert syscallnum to syscall name
func ToSyscallName(sysnum uint) (string, error) {
	if errInfo != nil {
		return "", errInfo
	}
	n, ok := info.SyscallNumbers[int(sysnum)]
	if !ok {
		return "", fmt.Errorf("syscall no %d does not exits", sysnum)
	}
	return n, nil
}
