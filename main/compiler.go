//go:build linux && go1.15
// +build linux,go1.15

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
	"unicode"

	"github.com/Sxu-Online-Judge/judge/model"
)

func main() {
	compiler("C")
}

func compiler(lang string) {
	script := model.Scripts[lang]

	var stdout, stderr bytes.Buffer

	cmdList := splitCmd(script.CompileScript)
	cmd := exec.Command(cmdList[0], cmdList[1:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = script.BaseDir

	time.AfterFunc(time.Duration(script.TimeOut)*time.Millisecond, func() {
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	})

	if err := cmd.Run(); err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("stderr: %s, err: %s\n", stderr.String(), err.Error()))
		return
	}

	_, _ = os.Stdout.WriteString("Compile OK\n")
}

func splitCmd(s string) (res []string) {
	var buf bytes.Buffer
	insideQuotes := false
	for _, r := range s {
		switch {
		case unicode.IsSpace(r) && !insideQuotes:
			if buf.Len() > 0 {
				res = append(res, buf.String())
				buf.Reset()
			}
		case r == '"' || r == '\'':
			if insideQuotes {
				res = append(res, buf.String())
				buf.Reset()
				insideQuotes = false
				continue
			}
			insideQuotes = true
		default:
			buf.WriteRune(r)
		}
	}
	if buf.Len() > 0 {
		res = append(res, buf.String())
	}
	return
}
