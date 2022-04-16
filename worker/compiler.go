package worker

import (
	"strings"

	"github.com/SXUOJ/judge/config"
	"github.com/SXUOJ/judge/lang"
	"github.com/SXUOJ/judge/pkg/seccomp"
	"github.com/SXUOJ/judge/runner"
	"github.com/SXUOJ/judge/sandbox"
)

type Compiler struct {
	realTimeLimit uint64
	r             *sandbox.Runner
}

func NewCompiler(worker *Worker, lang lang.Lang) (*Compiler, error) {
	rlimits, limit := parseLimit(lang.CompileCpuTimeLimit(), lang.CompileRealTimeLimit(), 0, 0, lang.CompileMemoryLimit())

	defaultAction, allow, trace, h := config.GetConf(strings.Join([]string{worker.Type, "compile"}, "-"), worker.AllowProc)
	seccompBuilder := seccomp.Builder{
		Default: defaultAction,
		Allow:   allow,
		Trace:   trace,
	}
	filter, _ := seccompBuilder.Build()

	return &Compiler{
		realTimeLimit: lang.CompileRealTimeLimit(),
		r: &sandbox.Runner{
			Args: lang.CompileArgs(),
			Env:  []string{pathEnv},
			// ExecFile:    execFile,
			// Files:       fds,
			// WorkDir:     worker.WorkDir,
			Seccomp:     filter,
			RLimits:     rlimits.PrepareRLimit(),
			Limit:       limit,
			ShowDetails: config.ShowDetails,
			Unsafe:      config.UnSafe,
			Handler:     h,
		},
	}, nil
}

func (Compiler *Compiler) Run() (*CompileResult, error) {
	res, err := run(Compiler.r, Compiler.realTimeLimit)
	if err != nil {
		return &CompileResult{
			Status: StatusCE,
			Error:  err.Error(),
		}, nil
	}
	if res.Status != runner.StatusNormal {
		return &CompileResult{
			Status: StatusCE,
		}, nil
	}
	return nil, nil
}
