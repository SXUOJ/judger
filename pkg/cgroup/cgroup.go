//Package cgroup provices builder to create cgroup

package cgroup

type Cgroup interface {
	SetCPUQuota() error
	SetCPUSet() error
	SetMemoryLimit() error
	SetProcLimit() error

	AddProc() error

	MemoryUsage() (uint64, error)
	MemoryMaxUsage() (uint64, error)
	CPUUsage() (uint64, error)

	Destroy() error
}
