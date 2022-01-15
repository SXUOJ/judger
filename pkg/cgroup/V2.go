package cgroup

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
)

type CgroupV2 struct {
	path string
}

func (c *CgroupV2) SetCPUSet(content string) error {
	return c.WriteFile("cpuset.cpus", []byte(content))
}

func (c *CgroupV2) SetCPUQuota(max, period uint64) error {
	content := strconv.FormatUint(max, 10) + " " + strconv.FormatUint(period, 10)
	return c.WriteFile("cpu.max", []byte(content))
}

func (c *CgroupV2) SetMemoryLimit(limit uint64) error {
	return c.WriteUint("memory.max", limit)
}

func (c *CgroupV2) SetProcLimit(limit uint64) error {
	return c.WriteUint("pids.max", limit)
}

func (c *CgroupV2) AddProc(pid uint64) error {
	return c.WriteUint("cgroup.procs", pid)
}

func (c *CgroupV2) CPUUsage() (uint64, error) {
	content, err := c.ReadFile("cpu.stat")
	if err != nil {
		return 0, err
	}

	temp := bufio.NewScanner(bytes.NewReader(content))
	for temp.Scan() {
		v := strings.Fields(temp.Text())
		if len(v) == 2 && v[0] == "usage_usec" {
			vv, err := strconv.Atoi(v[1])
			if err != nil {
				return 0, err
			}
			return uint64(vv) * 1000, nil
		}
	}
	return 0, os.ErrNotExist
}

func (c *CgroupV2) MemoryUsage() (uint64, error) {
	return c.ReadUint("memory.current")
}

//TODO:memory status
func (c *CgroupV2) MemoryMaxUsage() (uint64, error) {
	content, err := c.ReadFile("memory.stat")
	if err != nil {
		return 0, err
	}

	fmt.Printf("%v", content)
	return 0, nil
}

func (c *CgroupV2) Destroy() error {
	return remove(c.path)
}

func (c *CgroupV2) WriteUint(filename string, num uint64) error {
	return c.WriteFile(filename, []byte(strconv.FormatUint(num, 10)))
}

func (c *CgroupV2) WriteFile(filename string, content []byte) error {
	if c == nil || c.path == "" {
		return ErrNotExistence
	}
	p := path.Join(c.path, filename)
	return writeFile(p, content)
}

func (c *CgroupV2) ReadUint(filename string) (uint64, error) {
	content, err := c.ReadFile(filename)
	if err != nil {
		return 0, err
	}

	num, err := strconv.ParseUint(strings.TrimSpace(string(content)), 10, 64)
	if err != nil {
		return 0, err
	}
	return num, nil
}

func (c *CgroupV2) ReadFile(filename string) ([]byte, error) {
	if c == nil || filename == "" {
		return nil, ErrNotExistence
	}
	path := path.Join(c.path, filename)
	return readFile(path)
}
