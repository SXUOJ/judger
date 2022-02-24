package judge

import (
	"os"
	"path/filepath"

	"github.com/Sxu-Online-Judge/judge/model"
)

const (
	baseDir = "/home/ther/sandbox/tmp"
)

type Init interface {
	Init() model.JudgeStatus
}

type Preparer struct {
	codeType   string
	sourceCode string

	codeDir string
}

func (p *Preparer) Init() model.JudgeStatus {
	if err := os.MkdirAll(p.codeDir, os.ModePerm); err != nil {
		return model.SystemError
	}

	name := getNameByType(p.codeType)
	sourceCodeFile := filepath.Join(p.codeDir, name)
	out, err := os.OpenFile(sourceCodeFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0)
	defer out.Close()

	_, err = out.WriteString(p.sourceCode)
	if err != nil {
		return model.SystemError
	}

	return model.Normal
}

func getNameByType(langType string) string {
	switch langType {
	case "C":
		return "Main.c"
	case "Cpp":
		return "Main.cpp"
	case "Go":
		return "Main.go"
	case "Python2", "Python3":
		return "Main.py"
	}
	return ""
}
