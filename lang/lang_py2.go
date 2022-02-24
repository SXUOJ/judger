package lang

import (
	"fmt"
	"os/exec"
)

type Python2 struct {
	Dir string
}

func newPython2(dir string) *Python2 {
	return &Python2{
		Dir: dir,
	}
}

func (p *Python2) NeedCompile() bool {
	return false
}

func (p *Python2) Compile() *exec.Cmd {
	return nil
}

func (p *Python2) Run() *exec.Cmd {
	return exec.Command(
		"/usr/bin/python",
		fmt.Sprintf("%v/Main.py", p.Dir),
	)
}
