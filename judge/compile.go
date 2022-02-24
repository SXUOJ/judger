package judge

import (
	"bytes"
	"fmt"

	"github.com/Sxu-Online-Judge/judge/lang"
	"github.com/Sxu-Online-Judge/judge/model"
)

type CompileResult struct {
	status           model.JudgeStatus
	compileErrorInfo string
}

type Compiler struct {
	codeType string
	codeDir  string
}

func NewCompile(codeType, codeDir string) *Compiler {
	return &Compiler{
		codeType: codeType,
		codeDir:  codeDir,
	}
}

func (c *Compiler) Compile() (cr CompileResult) {
	lang, err := lang.NewLang(c.codeType, c.codeDir)
	if err != nil {
		cr.status = model.StatusCE
		cr.compileErrorInfo = err.Error()
		return
	}
	if !lang.NeedCompile() {
		cr.status = model.Normal
		return
	}

	compileCmd := lang.Compile()
	errInf := bytes.NewBuffer(nil)
	compileCmd.Stderr = errInf

	if err := compileCmd.Run(); err != nil {
		fmt.Println("compile error")
		cr.status = model.StatusCE
		cr.compileErrorInfo = errInf.String()
		return
	}

	cr.status = model.Normal
	return
}
