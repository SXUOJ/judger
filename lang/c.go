package lang

type c struct {
	compile
	runArgs
}

func newC(sourcePath, binaryPath string) *c {
	return &c{
		compile: compile{
			args: []string{
				"/usr/bin/gcc",
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
		runArgs: []string{""},
	}
}

func (c *c) NeedCompile() bool {
	return true
}

func (c *c) CompileArgs() []string {
	return c.compile.args
}

func (c *c) CompileRealTimeLimit() uint64 {
	return c.compile.realTimeLimit
}

func (c *c) CompileCpuTimeLimit() uint64 {
	return c.compile.cpuTimeLimit
}

func (c *c) CompileMemoryLimit() uint64 {
	return c.compile.memoryLimit
}

func (c *c) RunArgs() []string {
	return c.runArgs
}
