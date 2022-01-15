package cgroup

import (
	"testing"
	"time"
)

func TestNewBuilder(t *testing.T) {
	cgroupBuilder := NewBuilder().AddType("V2").AddCPU().AddCPUAcct().AddCPUSet().AddCPU().AddMemory().AddPids()

	cg, err := cgroupBuilder.Build("test")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		cg.Destroy()
	})

	if err := cg.SetCPUQuota(50000, 100000); err != nil {
		t.Fatal(err)
	}

	if err := cg.SetCPUSet("6"); err != nil {
		t.Fatal(err)
	}

	if err := cg.SetMemoryLimit(256); err != nil {
		t.Fatal(err)
	}

	if err := cg.SetProcLimit(10); err != nil {
		t.Fatal(err)
	}

	if err := cg.AddProc(uint64(138898)); err != nil {
		t.Fatal(err)
	}

	if _, err := cg.CPUUsage(); err != nil {
		t.Fatal(err)
	}

	if _, err := cg.MemoryUsage(); err != nil {
		t.Fatal(err)
	}

	if _, err := cg.MemoryMaxUsage(); err != nil {
		t.Fatal(err)
	}
	time.Sleep(5000 * time.Millisecond)
}
