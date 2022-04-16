package worker

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/SXUOJ/judge/lang"
	"github.com/SXUOJ/judge/runner"
	"github.com/SXUOJ/judge/sandbox"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Worker struct {
	ProblemID string
	SubmitID  string

	FileName string
	Type     string
	WorkDir  string

	AllowProc     bool
	TimeLimit     uint64
	RealTimeLimit uint64
	MemoryLimit   uint64
	OutputLimit   uint64
	StackLimit    uint64

	compiler Compiler
	runner   Runner
}

func (worker *Worker) Run(c *gin.Context) {
	worker.WorkDir = filepath.Join(runDir, worker.SubmitID)
	var (
		sourcePath = filepath.Join(worker.WorkDir, worker.FileName)
		binaryPath = filepath.Join(worker.WorkDir, worker.SubmitID)
	)

	lang, err := lang.NewLang(worker.Type, sourcePath, binaryPath)
	if err != nil {
		msg := "New lang failed"
		errReturn(c, msg)
		return
	}

	if lang.NeedCompile() {
		compiler, err := NewCompiler(worker, lang, c)
		if err != nil {
			msg := "New compiler failed"
			errReturn(c, msg)
			return
		}
		if err := compiler.Run(); err != nil {
			logrus.Error(err.Error())
			return
		}
	}

	run, err := NewRunner(worker, lang, c)
	if err != nil {
		msg := "New runner failed"
		errReturn(c, msg)
		return
	}
	run.Run()
	return
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

func errReturn(c *gin.Context, msg string) {
	logrus.Error(msg)
	c.JSON(http.StatusOK, gin.H{
		"msg": msg,
	})
}
