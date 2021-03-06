package sandbox

import (
	"context"

	"github.com/SXUOJ/judge/pkg/forkexec"
	"github.com/SXUOJ/judge/pkg/rlimit"
	"github.com/SXUOJ/judge/pkg/seccomp"
	"github.com/SXUOJ/judge/ptrace"
	"github.com/SXUOJ/judge/runner"
)

type Runner struct {
	Args     []string
	Env      []string
	WorkDir  string
	ExecFile uintptr
	Files    []uintptr
	RLimits  []rlimit.RLimit
	Seccomp  seccomp.Filter
	SyncFunc func(pid int) error

	Handler             Handler
	ShowDetails, Unsafe bool

	Limit runner.Limit
}

type Handler interface {
	CheckSyscall(string) ptrace.TraceAction
}

func (r *Runner) Run(c context.Context) runner.Result {
	fork := &forkexec.Runner{
		Args:     r.Args,
		Env:      r.Env,
		ExecFile: r.ExecFile,
		RLimits:  r.RLimits,
		Files:    r.Files,
		WorkDir:  r.WorkDir,
		Seccomp:  r.Seccomp.SockFprog(),
		Ptrace:   true,
		SyncFunc: r.SyncFunc,
	}

	th := &tracerHandler{
		ShowDetails: r.ShowDetails,
		Unsafe:      r.Unsafe,
		Handler:     r.Handler,
	}

	tracer := ptrace.Ptracer{
		Handler: th,
		Runner:  fork,
		Limit:   r.Limit,
	}
	return tracer.Trace(c)
}
