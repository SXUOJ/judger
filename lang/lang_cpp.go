package lang

import (
	"fmt"
	"os/exec"
)

type Cpp struct {
	Dir string
}

func newCpp(dir string) *C {
	return &C{
		Dir: dir,
	}
}

func (c *Cpp) NeedCompile() bool {
	return true
}

func (c *Cpp) Compile() *exec.Cmd {
	return exec.Command(
		"/usr/bin/g++",
		"-o",
		fmt.Sprintf("%v/Main", c.Dir),
		fmt.Sprintf("%v/Main.cpp", c.Dir),
		"-fmax-errors=3",
		"-std=c11",
		"-lm",
		"-w",
		"-O2",
		"-DONLINE_JUDGE",
	)
}

func (c *Cpp) Run() *exec.Cmd {
	return exec.Command(
		fmt.Sprintf("%v/Main", c.Dir),
	)
}
