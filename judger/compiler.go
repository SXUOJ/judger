package judger

import (
	"strings"
	"time"

	"github.com/SXUOJ/judge/config"
	"github.com/SXUOJ/judge/lang"
	"github.com/SXUOJ/judge/pkg/seccomp"
	"github.com/SXUOJ/judge/runner"
	"github.com/SXUOJ/judge/sandbox"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type compileResult runner.Result

type compiler struct {
	submitID      string
	realTimeLimit uint64
	r             *sandbox.Runner
}

func newCompiler(submitID, codeType, workDir string, allowProc bool, lang lang.Lang, c *gin.Context) *compiler {
	rlimits, limit := parseLimit(lang.CompileCpuTimeLimit(), lang.CompileRealTimeLimit(), 0, 0, lang.CompileMemoryLimit())

	defaultAction, allow, trace, h := config.GetConf(strings.Join([]string{codeType, "compile"}, "-"), allowProc)
	seccompBuilder := seccomp.Builder{
		Default: defaultAction,
		Allow:   allow,
		Trace:   trace,
	}
	filter, _ := seccompBuilder.Build()

	return &compiler{
		submitID:      submitID,
		realTimeLimit: lang.CompileRealTimeLimit(),
		r: &sandbox.Runner{
			Args:        lang.CompileArgs(),
			Env:         []string{pathEnv},
			WorkDir:     workDir,
			Seccomp:     filter,
			RLimits:     rlimits.PrepareRLimit(),
			Limit:       limit,
			ShowDetails: config.ShowDetails,
			Unsafe:      config.UnSafe,
			Handler:     h,
		},
	}
}

func (compile compiler) Start(resChan chan interface{}) {
	logrus.Info("Start Compile: ", compile.submitID)
	defer logrus.Info("End compile: ", compile.submitID)
	res, err := run(compile.r, compile.realTimeLimit)

	if err != nil || res.Status != runner.StatusNormal {
		result := compileResult{
			Status:      runner.StatusCompileError,
			SetUpTime:   res.SetUpTime,
			RunningTime: res.RunningTime / time.Millisecond,
			Time:        res.Time / time.Millisecond,
			Memory:      res.Memory >> 20,
			ExitCode:    res.ExitCode,
			Error:       res.Error,
		}
		resChan <- result
		return
	}
	resChan <- nil
}
