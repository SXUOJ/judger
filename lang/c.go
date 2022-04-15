package lang

type C struct {
	compile
	run
}

func newC(sourcePath, binaryPath string) *C {
	return &C{
		compile: compile{
			bin: "/usr/bin/gcc",
			args: []string{
				"-o",
				binaryPath,
				sourcePath,
				"-fmax-errors=3",
				"-std=c11",
				"-lm",
				"-w",
				"-O2",
				"-DONLINE_JUDGE",
				"",
			},
			realTimeLimit: 5,
			cpuTimeLimit:  3,
			memoryLimit:   128 * 1024 * 1024,
		},
		run: run{
			bin:  binaryPath,
			args: []string{""},
		},
	}
}

func (c *C) IsCompile() bool {
	return true
}

func (c *C) CompileBin() string {
	return c.compile.bin
}

func (c *C) CompileArgs() []string {
	return c.compile.args
}

func (c *C) CompileRealTimeLimit() uint64 {
	return c.compile.realTimeLimit
}

func (c *C) CompileCpuTimeLimit() uint64 {
	return c.compile.cpuTimeLimit
}

func (c *C) CompileMemoryLimit() uint64 {
	return c.compile.memoryLimit
}

func (c *C) RunBin() string {
	return c.run.bin
}

func (c *C) RunArgs() []string {
	return c.run.args
}
