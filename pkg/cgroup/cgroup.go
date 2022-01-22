//Package cgroup provices builder to create cgroup

package cgroup

import "errors"

type Cgroup interface {
	SetCPUSet(string) error
	SetCPUQuota(uint64) error
	SetMemoryLimit(uint64) error
	SetProcLimit(uint64) error

	AddProc(uint64) error

	CPUUsage() (uint64, error)
	MemoryUsage() (uint64, error)
	MemoryMaxUsage() (uint64, error)

	Destroy() error
}

type ResourceConfig struct {
	MemoryLimit string
	CpuSet      string
	CpuQuota    string
}

const (
	basePathV2  = "/sys/fs/cgroup"
	cgroupProcs = "cgroup.procs"

	filePerm = 0644
	dirPerm  = 0755
)

var (
	ErrNotExistence = errors.New("This path not exist")
)
