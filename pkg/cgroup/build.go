package cgroup

import (
	"log"
	"os"
	"path"
)

type CgroupType int

const (
	_ = iota
	CgroupTypeV1
	CgroupTypeV2
)

type Builder struct {
	Type CgroupType

	CPU     bool
	CPUAcct bool
	CPUSet  bool
	Memory  bool
	Pids    bool
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (builder *Builder) AddType(cgt string) *Builder {
	switch cgt {
	case "V1", "v1":
		builder.Type = 1
		return builder
	case "V2", "v2":
		builder.Type = 2
		return builder
	default:
		log.Fatal("cgroup type error")
		return nil
	}
}
func (builder *Builder) AddCPU() *Builder     { builder.CPU = true; return builder }
func (builder *Builder) AddCPUAcct() *Builder { builder.CPUAcct = true; return builder }
func (builder *Builder) AddCPUSet() *Builder  { builder.CPUSet = true; return builder }
func (builder *Builder) AddMemory() *Builder  { builder.Memory = true; return builder }
func (builder *Builder) AddPids() *Builder    { builder.Pids = true; return builder }

func (builder *Builder) Build(name string) (Cgroup, error) {
	if builder.Type == CgroupTypeV1 {
		return builder.buildV1(name)
	} else {
		return builder.buildV2(name)
	}
}

func (builder *Builder) buildV1(name string) (Cgroup, error) { return nil, nil }

func (builder *Builder) buildV2(name string) (cg Cgroup, err error) {
	path := path.Join(basePathV2, name)
	defer func() {
		if err != nil {
			remove(path)
		}
	}()

	if err := os.Mkdir(path, 0755); err != nil {
		return nil, err
	}

	return &CgroupV2{path}, nil
}
