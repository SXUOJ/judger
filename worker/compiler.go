package worker

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/SXUOJ/judge/config"
	"github.com/SXUOJ/judge/lang"
	"github.com/SXUOJ/judge/pkg/seccomp"
	"github.com/SXUOJ/judge/runner"
	"github.com/SXUOJ/judge/sandbox"
	"github.com/gin-gonic/gin"
)

type CompileResult runner.Result

type Compiler struct {
	submit_id     string
	realTimeLimit uint64
	r             *sandbox.Runner
	c             *gin.Context
}

func NewCompiler(worker *Worker, lang lang.Lang, c *gin.Context) (*Compiler, error) {
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
			WorkDir:     worker.WorkDir,
			Seccomp:     filter,
			RLimits:     rlimits.PrepareRLimit(),
			Limit:       limit,
			ShowDetails: config.ShowDetails,
			Unsafe:      config.UnSafe,
			Handler:     h,
		},
		c: c,
	}, nil
}

func (compile *Compiler) Run() error {
	res, err := run(compile.r, compile.realTimeLimit)
	if err != nil || res.Status != runner.StatusNormal {
		compile.c.JSON(http.StatusOK, gin.H{
			"msg":       "compile error",
			"submit_id": compile.submit_id,
			"result": CompileResult{
				Status:      runner.StatusCompileError,
				SetUpTime:   res.SetUpTime,
				RunningTime: res.RunningTime / time.Millisecond,
				Time:        res.Time / time.Millisecond,
				Memory:      res.Memory >> 20,
				ExitCode:    res.ExitCode,
				Error:       res.Error,
			},
		})
		return errors.New("Compile Error")
	}
	return nil
}
