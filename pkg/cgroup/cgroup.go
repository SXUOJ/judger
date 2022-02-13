//Package cgroup provices builder to create cgroup

package cgroup

import "errors"

type Cgroup interface {
	Init(string, string) error

	Destroy() error
}

type ResourceConfig struct {
	MemoryLimit string
	CpuSet      string
	CpuQuota    string
}

const (
	basePathV2 = "/sys/fs/cgroup"

	cpuPrefixV1 = "/sys/fs/cgroup/cpu/"
	pidPrefixV1 = "/sys/fs/cgroup/pids/"
	memPrefixV1 = "/sys/fs/cgroup/memory/"

	filePerm = 0644
	dirPerm  = 0755
)

var (
	ErrNotExistence = errors.New("This path not exist")
)
