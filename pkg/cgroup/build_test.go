package cgroup

import (
	"log"
	"testing"
)

func TestNewBuilder(t *testing.T) {
	cgroupBuilder := NewBuilder().AddCPU().AddCPUAcct().AddCPUSet().AddCPU().AddMemory().AddPids()
	log.Println(cgroupBuilder)
}
