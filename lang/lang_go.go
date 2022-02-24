package lang

import (
	"fmt"
	"os/exec"
)

type Go struct {
	Dir string
}

func newGo(dir string) *Go {
	return &Go{
		Dir: dir,
	}
}

func (g *Go) NeedCompile() bool {
	return true
}

func (g *Go) Compile() *exec.Cmd {
	return exec.Command(
		"/usr/bin/go",
		"build",
		"-o",
		fmt.Sprintf("%v/Main", g.Dir),
		fmt.Sprintf("%v/Main.go", g.Dir),
	)
}

func (g *Go) Run() *exec.Cmd {
	return exec.Command(
		fmt.Sprintf("%v/Main", g.Dir),
	)
}
