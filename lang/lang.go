package lang

import "errors"

type Lang interface {
	IsCompile() bool

	CompileBin() string
	CompileArgs() []string
	CompileRealTimeLimit() uint64
	CompileCpuTimeLimit() uint64
	CompileMemoryLimit() uint64

	RunBin() string
	RunArgs() []string
}

type compile struct {
	bin           string
	args          []string
	realTimeLimit uint64
	cpuTimeLimit  uint64
	memoryLimit   uint64
}

type run struct {
	bin  string
	args []string
}

func NewLang(cType, sourcePath, binaryPath string) (Lang, error) {
	switch cType {
	case "C", "c":
		return newC(sourcePath, binaryPath), nil
	default:
		return nil, errors.New("No type")
	}
}
