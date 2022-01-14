package cgroup

type CgroupV1 struct{}

//TODO: V1 function
func (c *CgroupV1) SetCPUQuota()    {}
func (c *CgroupV1) SetCPUSet()      {}
func (c *CgroupV1) SetMemoryLimit() {}
func (c *CgroupV1) SetProcLimit()   {}

func (c *CgroupV1) AddProc() {}

func (c *CgroupV1) CPUUsage()       {}
func (c *CgroupV1) MemoryUsage()    {}
func (c *CgroupV1) MemoryMaxUsage() {}

func (c *CgroupV1) Destroy() {}
