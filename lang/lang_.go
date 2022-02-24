package lang

import (
	"errors"
	"os/exec"
)

var (
	ERROR_NOT_SUPPORT_LANGUAGE = errors.New("This language is not supported")
)

const (
	_ = iota
	LangC
	LangCpp
	LangJava
	LangPython2
	LangPython3
	LangGo
)

type Lang interface {
	NeedCompile() bool

	Compile() *exec.Cmd
	Run() *exec.Cmd
}

func NewLang(langType string, langDir string) (Lang, error) {
	switch langType {
	case "C":
		return newC(langDir), nil
	case "Cpp":
		return newCpp(langDir), nil
	case "Go":
		return newGo(langDir), nil
	case "Python2":
		return newPython2(langDir), nil
	case "Python3":
		return newPython3(langDir), nil
	default:
		return nil, ERROR_NOT_SUPPORT_LANGUAGE
	}
}
