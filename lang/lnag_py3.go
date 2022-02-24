package lang

import (
	"fmt"
	"os/exec"
)

type Python3 struct {
	Dir string
}

func newPython3(dir string) *Python3 {
	return &Python3{
		Dir: dir,
	}
}

func (p *Python3) NeedCompile() bool {
	return false
}

func (p *Python3) Compile() *exec.Cmd {
	return nil
}

func (p *Python3) Run() *exec.Cmd {
	return exec.Command(
		"/usr/bin/python3",
		fmt.Sprintf("%v/Main.py", p.Dir),
	)
}
