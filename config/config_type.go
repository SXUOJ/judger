package config

// ProgramConfig defines the extra config apply to program type
type ProgramConfig struct {
	Syscall SyscallConfig
}

// SyscallConfig defines extra syscallConfig apply to program type
type SyscallConfig struct {
	Allow []string
	Trace []string
}

// FileAccessConfig defines extra file access permission for the program type
type FileAccessConfig struct {
	ExtraRead, ExtraWrite, ExtraStat, ExtraBan []string
}
