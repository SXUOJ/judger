package lang

import "errors"

type Lang interface {
	NeedCompile() bool

	CompileArgs() []string
	CompileRealTimeLimit() uint64
	CompileCpuTimeLimit() uint64
	CompileMemoryLimit() uint64

	RunArgs() []string
}

type compile struct {
	args          []string
	realTimeLimit uint64
	cpuTimeLimit  uint64
	memoryLimit   uint64
}

type runArgs []string

func NewLang(cType, sourcePath, binaryPath string) (Lang, error) {
	switch cType {
	case "C", "c":
		return newC(sourcePath, binaryPath), nil
	default:
		return nil, errors.New("No type")
	}
}
