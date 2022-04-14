package forkexec

import (
	"os"
	"testing"

	"github.com/SXUOJ/judge/pkg/seccomp"
	"golang.org/x/sys/unix"
)

func TestFork(t *testing.T) {
	t.Parallel()
	f, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	if err := f.Chmod(0777); err != nil {
		t.Fatal(err)
	}

	echo, err := os.Open("/bin/echo")
	if err != nil {
		t.Fatal(err)
	}
	defer echo.Close()

	_, err = f.ReadFrom(echo)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	r := Runner{
		Args: []string{f.Name()},
	}
	_, err = r.Start()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSeccomp(t *testing.T) {
	var (
		args = []string{"../../test/resources/c/hello"}
	)

	fds := []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()}
	filter := getWriteFileter(true) //true or false

	r := &Runner{
		Args:    args,
		Env:     os.Environ(),
		Files:   fds,
		Ptrace:  true,
		Seccomp: filter.SockFprog(),
	}

	pid, err := r.Start()
	if err != nil {
		t.Fatal(err)
	}
	var ws unix.WaitStatus
	_, err = unix.Wait4(pid, &ws, 0, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func getWriteFileter(allow bool) *seccomp.Filter {
	var (
		trace = []string{}
	)
	if !allow {
		trace = append(trace, "write")
	}

	seccompBuilder := seccomp.Builder{
		Default: seccomp.ActionAllow,
		Trace:   trace,
	}
	filter, err := seccompBuilder.Build()
	if err != nil {
		panic(err)
	}
	return &filter
}
