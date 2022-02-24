package lang

import (
	"fmt"
	"os/exec"
)

type C struct {
	Dir string
}

func newC(dir string) *C {
	return &C{
		Dir: dir,
	}
}

func (c *C) NeedCompile() bool {
	return true
}

func (c *C) Compile() *exec.Cmd {
	return exec.Command(
		"/usr/bin/gcc",
		"-o",
		fmt.Sprintf("%v/Main", c.Dir),
		fmt.Sprintf("%v/Main.c", c.Dir),
		"-fmax-errors=3",
		"-std=c11",
		"-lm",
		"-w",
		"-O2",
		"-DONLINE_JUDGE",
	)
}

func (c *C) Run() *exec.Cmd {
	return exec.Command(
		fmt.Sprintf("%v/Main", c.Dir),
	)
}
