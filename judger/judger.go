package judger

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/SXUOJ/judge/lang"
	"github.com/SXUOJ/judge/runner"
	"github.com/SXUOJ/judge/sandbox"
	"github.com/SXUOJ/judge/worker"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Judger struct {
	SubmitID string

	FileName string
	Type     string
	WorkDir  string

	AllowProc bool
	Slimit    Limit
}

type Limit struct {
	TimeLimit     uint64
	RealTimeLimit uint64
	OutputLimit   uint64
	StackLimit    uint64
	MemoryLimit   uint64
}

func NewJudger() *Judger {
	return nil
}

func (judger *Judger) Run(c *gin.Context) {
	var (
		sourcePath = filepath.Join(judger.WorkDir, judger.FileName)
		binaryPath = filepath.Join(judger.WorkDir, judger.SubmitID)
	)

	lang, err := lang.NewLang(judger.Type, sourcePath, binaryPath)
	if err != nil {
		msg := "New lang failed"
		logrus.Info(msg)
		return
	}

	compiler := newCompiler(judger.SubmitID, judger.Type, judger.WorkDir, judger.AllowProc, lang, c)
	thisRunner := newRunner(judger.SubmitID, judger.Type, judger.WorkDir, judger.AllowProc, judger.Slimit, lang, c)

	w := &worker.Worker{
		C:        c,
		Done:     make(chan bool),
		Lang:     lang,
		Compiler: compiler,
		Runner:   thisRunner,
	}
	worker.WorkPool.Put(w)
	w.Wait()
}

func run(r *sandbox.Runner, realTimeLimit uint64) (*runner.Result, error) {
	var (
		rt runner.Result
	)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	// Run tracer
	sTime := time.Now()
	c, cancel := context.WithTimeout(context.Background(), time.Duration(int64(realTimeLimit)*int64(time.Second)))
	defer cancel()

	s := make(chan runner.Result, 1)
	go func() {
		s <- r.Run(c)
	}()
	rTime := time.Now()

	select {
	case <-sig:
		cancel()
		rt = <-s
		rt.Status = runner.StatusSystemError

	case rt = <-s:
	}
	eTime := time.Now()

	if rt.SetUpTime == 0 {
		rt.SetUpTime = rTime.Sub(sTime)
		rt.RunningTime = eTime.Sub(rTime)
	}

	return &rt, nil
}
