package cgroup

import (
	"log"
	"path"

	"github.com/sirupsen/logrus"
)

type CgroupV2 struct {
	name string
	path string
}

func (c *CgroupV2) Init(pid, memory string) error {
	logrus.Info("init cgroup start")

	if err := c.setCPUQuota("10000"); err != nil {
		log.Println("SetCPUQuota failed")
		return err
	}

	if err := c.setMemoryLimit(memory); err != nil {
		logrus.Error("SetMemoryLimit failed")
		return err
	}

	if err := c.addProc(pid); err != nil {
		logrus.Error("AddProc failed")
		return err
	}

	logrus.Info("init cgroup done")
	return nil
}

func (c *CgroupV2) setCPUQuota(period string) error {
	return c.writeFile("cpu.max", []byte(period))
}

func (c *CgroupV2) setMemoryLimit(limit string) error {
	return c.writeFile("memory.max", []byte(limit))
}

func (c *CgroupV2) addProc(pid string) error {
	return c.writeFile("cgroup.procs", []byte(pid))
}

func (c *CgroupV2) Destroy() error {
	return remove(c.path)
}

func (c *CgroupV2) writeFile(filename string, content []byte) error {
	if c == nil || c.path == "" {
		return ErrNotExistence
	}
	p := path.Join(c.path, filename)
	return writeFile(p, content)
}
