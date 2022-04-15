package config

import "github.com/SXUOJ/judge/pkg/seccomp"

// This file includes configs for the run program settings

var (
	// default allowed safe syscalls
	defaultSyscallAllows = []string{}

	// default syscalls to trace
	defaultSyscallTraces = []string{}

	// process related syscall if allowProc enabled
	defaultProcSyscalls = []string{"clone", "fork", "vfork", "nanosleep", "execve"}

	// config for different type of program
	// workpath and arg0 have additional read / stat permission
	runptraceConfig = map[string]SyscallConfig{
		"default": {
			DefaultAction: seccomp.ActionAllow,
			Trace:         []string{
				// "write",
			},
			Allow: []string{},
		},
		"C-compile": {
			DefaultAction: seccomp.ActionAllow,
			Trace:         []string{},
			Allow:         []string{},
		},
	}
)
