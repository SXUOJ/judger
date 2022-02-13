package cgroup

import (
	"path"

	"github.com/sirupsen/logrus"
)

type CgroupV2 struct {
	name string
	path string
}

var (
	cpuFileName = "cpu.max"
	memFileName = "memory.max"
	pidFileName = "cgroup.procs"
)

func (c *CgroupV2) Init(pid, memory string) error {
	logrus.Info("Init cgroup start")

	if err := c.setCPUQuota("10000"); err != nil {
		logrus.Error("func setCPUQuota() failed")
		return err
	}

	if err := c.setMemoryLimit(memory); err != nil {
		logrus.Error("func setMemoryLimit() failed")
		return err
	}

	if err := c.addProc(pid); err != nil {
		logrus.Error("func addProc() failed")
		return err
	}

	logrus.Info("Init cgroup done")
	return nil
}

func (c *CgroupV2) Destroy() error {
	return remove(c.path)
}

func (c *CgroupV2) setCPUQuota(period string) error {
	return c.writeFile(cpuFileName, []byte(period))
}

func (c *CgroupV2) setMemoryLimit(limit string) error {
	return c.writeFile(memFileName, []byte(limit))
}

func (c *CgroupV2) addProc(pid string) error {
	return c.writeFile(pidFileName, []byte(pid))
}

func (c *CgroupV2) writeFile(filename string, content []byte) error {
	if c == nil || c.path == "" {
		return ErrNotExistence
	}
	p := path.Join(c.path, filename)
	return writeFile(p, content)
}
