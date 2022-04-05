// Package seccomp provides a generated filter format for seccomp filter
package seccomp

import (
	"syscall"
)

// Filter is the BPF seccomp filter value
type Filter []syscall.SockFilter

// SockFprog converts Filter to SockFprog for seccomp syscall
func (f Filter) SockFprog() *syscall.SockFprog {
	b := []syscall.SockFilter(f)
	return &syscall.SockFprog{
		Len:    uint16(len(b)),
		Filter: &b[0],
	}
}
