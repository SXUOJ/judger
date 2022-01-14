package cgroup

type CgroupV2 struct{}

//TODO: V2 function
func (c *CgroupV2) SetCPUQuota()    {}
func (c *CgroupV2) SetCPUSet()      {}
func (c *CgroupV2) SetMemoryLimit() {}
func (c *CgroupV2) SetProcLimit()   {}

func (c *CgroupV2) AddProc() {}

func (c *CgroupV2) CPUUsage()       {}
func (c *CgroupV2) MemoryUsage()    {}
func (c *CgroupV2) MemoryMaxUsage() {}

func (c *CgroupV2) Destroy() {}
