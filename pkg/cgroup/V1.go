package cgroup

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type CgroupV1 struct {
	name string
}

func (c *CgroupV1) Init(pid, memory string) error {
	logrus.Info("Init cgroup start")

	if err := cpuCgroup(pid, c.name); err != nil {
		logrus.Error("func cpuCgroup() failed")
		return err
	}

	if err := pidCgroup(pid, c.name); err != nil {
		logrus.Error("func pidCgroup() failed")
		return err
	}

	if err := memoryCgroup(pid, c.name, memory); err != nil {
		logrus.Error("func memoryCgroup() failed")
		return err
	}

	logrus.Info("Init cgroup done")
	return nil
}

func (c *CgroupV1) Destroy() error {
	dirs := []string{
		filepath.Join(cpuPrefixV1, c.name),
		filepath.Join(pidPrefixV1, c.name),
		filepath.Join(memPrefixV1, c.name),
	}

	var err error
	for _, dir := range dirs {
		if err = remove(dir); err != nil {
			logrus.Error("Remove cgroup failed")
		}
	}
	return err
}

func cpuCgroup(pid, containerID string) error {
	cgCPUPath := filepath.Join(cpuPrefixV1, containerID)
	mapping := map[string]string{
		"tasks":            pid,
		"cpu.cfs_quota_us": "10000",
	}

	for key, value := range mapping {
		path := filepath.Join(cgCPUPath, key)
		if err := ioutil.WriteFile(path, []byte(value), 0644); err != nil {
			return err
		}
	}
	return nil
}

func pidCgroup(pid, containerID string) error {
	cgPidPath := filepath.Join(pidPrefixV1, containerID)
	mapping := map[string]string{
		"cgroup.procs": pid,
		"pids.max":     "64",
	}

	for key, value := range mapping {
		path := filepath.Join(cgPidPath, key)
		if err := ioutil.WriteFile(path, []byte(value), 0644); err != nil {
			return err
		}
	}
	return nil
}

func memoryCgroup(pid, containerID, memory string) error {
	cgMemoryPath := filepath.Join(memPrefixV1, containerID)
	mapping := map[string]string{
		"memory.kmem.limit_in_bytes": "64m",
		"tasks":                      pid,
		"memory.limit_in_bytes":      fmt.Sprintf("%sm", memory),
	}

	for key, value := range mapping {
		path := filepath.Join(cgMemoryPath, key)
		if err := ioutil.WriteFile(path, []byte(value), 0644); err != nil {
			return err
		}
	}
	return nil
}
