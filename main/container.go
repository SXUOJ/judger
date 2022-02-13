//go:build linux && go1.15
// +build linux,go1.15

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/Sxu-Online-Judge/judge/model"
	"github.com/Sxu-Online-Judge/judge/pkg/cgroup"
	"github.com/Sxu-Online-Judge/judge/pkg/namespace"
	"github.com/docker/docker/pkg/reexec"
	uuid "github.com/gofrs/uuid"
)

func init() {
	// register "justiceInit" => justiceInit() every time
	reexec.Register("judgeInit", judgeInit)

	if reexec.Init() {
		os.Exit(0)
	}
}

func judgeInit() {
	basedir := os.Args[1]
	input := os.Args[2]
	expected := os.Args[3]
	timeout, _ := strconv.ParseInt(os.Args[4], 10, 32)

	r := new(model.Result)
	if err := namespace.Init(basedir); err != nil {
		result, _ := json.Marshal(r.GetRuntimeErrorTaskResult())
		_, _ = os.Stdout.Write(result)
		os.Exit(0)
	}

	var o, e bytes.Buffer
	cmd := exec.Command("/Main")
	cmd.Stdin = strings.NewReader(input)
	cmd.Stdout = &o
	cmd.Stderr = &e
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	cmd.Env = []string{"PS1=[judge] # "}

	time.AfterFunc(time.Duration(timeout)*time.Millisecond, func() {
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	})

	startTime := time.Now().UnixNano() / 1e6
	if err := cmd.Run(); err != nil {
		result, _ := json.Marshal(r.GetRuntimeErrorTaskResult())
		_, _ = os.Stdout.Write(result)
		_, _ = os.Stderr.WriteString(fmt.Sprintf("err: %s\n", err.Error()))
		return
	}
	endTime := time.Now().UnixNano() / 1e6

	if e.Len() > 0 {
		result, _ := json.Marshal(r.GetRuntimeErrorTaskResult())
		_, _ = os.Stdout.Write(result)
		_, _ = os.Stderr.WriteString(fmt.Sprintf("stderr: %s\n", e.String()))
		return
	}

	output := strings.TrimSpace(o.String())
	if output == expected {
		// ms, MB
		timeCost, memoryCost := endTime-startTime, cmd.ProcessState.SysUsage().(*syscall.Rusage).Maxrss/1024
		// timeCost value 0 will be omitted
		if timeCost == 0 {
			timeCost = 1
		}

		result, _ := json.Marshal(r.GetAcceptedTaskResult(timeCost, memoryCost))
		_, _ = os.Stdout.Write(result)
	} else {
		result, _ := json.Marshal(r.GetWrongAnswerTaskResult(input, output, expected))
		_, _ = os.Stdout.Write(result)
	}

	_, _ = os.Stderr.WriteString(fmt.Sprintf("output: %s | expected: %s\n", output, expected))
}

// logs will be printed to os.Stderr
func main() {
	basedir := flag.String("basedir", "/tmp", "basedir of tmp C binary")
	input := flag.String("input", "<input>", "test case input")
	expected := flag.String("expected", "<expected>", "test case expected")
	timeout := flag.String("timeout", "2000", "timeout in milliseconds")
	memory := flag.String("memory", "256", "memory limitation in MB")
	flag.Parse()

	result := new(model.Result)
	u, _ := uuid.NewV4()

	builder := cgroup.NewBuilder().AddType("V1")
	cg, _ := builder.Build(u.String())

	if err := cg.Init(strconv.Itoa(os.Getpid()), *memory); err != nil {
		result, _ := json.Marshal(result.GetRuntimeErrorTaskResult())
		_, _ = os.Stdout.Write(result)
		os.Exit(0)
	}

	cmd := reexec.Command("judgeInit", *basedir, *input, *expected, *timeout, *memory)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWUSER,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
	}

	if err := cmd.Run(); err != nil {
		result, _ := json.Marshal(result.GetRuntimeErrorTaskResult())
		_, _ = os.Stderr.WriteString(fmt.Sprintf("%s\n", err.Error()))
		_, _ = os.Stdout.Write(result)
	}

	os.Exit(0)
}
