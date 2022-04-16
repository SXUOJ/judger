package config

import "github.com/SXUOJ/judge/pkg/seccomp"

// This file includes configs for the run program settings

var (
	// default allowed safe syscalls
	defaultSyscallAllows = []string{
		"read",
		"write",
		"readv",
		"writev",
		"close",
		"fstat",
		"lseek",
		"dup",
		"dup2",
		"dup3",
		"ioctl",
		"fcntl",
		"fadvise64",
		"pread64",
		"pwrite64",

		// memory action
		"mmap",
		"mprotect",
		"munmap",
		"brk",
		"mremap",
		"msync",
		"mincore",
		"madvise",

		// signal action
		"rt_sigaction",
		"rt_sigprocmask",
		"rt_sigreturn",
		"rt_sigpending",
		"sigaltstack",

		// get current work dir
		"getcwd",

		// process exit
		"exit",
		"exit_group",

		// others
		"arch_prctl",

		"gettimeofday",
		"getrlimit",
		"getrusage",
		"times",
		"time",
		"clock_gettime",

		"restart_syscall",
	}

	// default syscalls to trace
	defaultSyscallTraces = []string{}

	// process related syscall if allowProc enabled
	defaultProcSyscalls = []string{"clone", "fork", "vfork", "nanosleep", "execve"}

	// config for different type of program
	// workpath and arg0 have additional read / stat permission
	runptraceConfig = map[string]SyscallConfig{
		"default": {
			DefaultAction: seccomp.ActionAllow,
			Trace:         []string{},
			Allow:         []string{},
		},
		"C-compile": {
			DefaultAction: seccomp.ActionAllow,
			Trace:         []string{},
			Allow:         []string{},
		},
		"C-run": {
			DefaultAction: seccomp.ActionAllow,
			Trace:         []string{},
			Allow:         []string{},
		},
	}
)
