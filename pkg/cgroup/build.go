package cgroup

import (
	"os"
	"path"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type CgroupType int

const (
	_ = iota
	CgroupTypeV1
	CgroupTypeV2
)

type Builder struct {
	Type CgroupType
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
		logrus.Error("cgroup type error")
		return nil
	}
}

func (builder *Builder) Build(name string) (Cgroup, error) {
	if builder.Type == CgroupTypeV1 {
		return builder.buildV1(name)
	} else {
		return builder.buildV2(name)
	}
}

func (builder *Builder) buildV1(name string) (cg Cgroup, err error) {
	dirs := []string{
		filepath.Join(cpuPrefixV1, name),
		filepath.Join(pidPrefixV1, name),
		filepath.Join(memPrefixV1, name),
	}
	defer func() {
		if err != nil {
			for _, dir := range dirs {
				remove(dir)
			}
		}
	}()

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, err
		}
	}
	return &CgroupV1{
		name: name,
	}, nil
}

func (builder *Builder) buildV2(name string) (cg Cgroup, err error) {
	path := path.Join(basePathV2, name)
	defer func() {
		if err != nil {
			remove(path)
		}
	}()

	// if err := os.Mkdir(path, dirPerm); err != nil {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return nil, err
	}

	return &CgroupV2{
		name: name,
		path: path,
	}, nil
}
