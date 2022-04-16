package worker

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/SXUOJ/judge/lang"
	"github.com/SXUOJ/judge/runner"
	"github.com/SXUOJ/judge/sandbox"
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
}

/*
 *|work_dir
 *|--submit_id
 *|--|--submit_id.c
 *|--|--submit_id
 *|--|--1.txt
 *|--|--2.txt
 *
 *|sample
 *|--problem_id
 *|--|--1.in
 *|--|--1.out
 *|--|--2.in
 *|--|--2.out
 */

func (worker *Worker) Run() (Results, error) {
	worker.WorkDir = filepath.Join(runDir, worker.SubmitID)
	var (
		sourcePath = filepath.Join(worker.WorkDir, worker.FileName)
		binaryPath = filepath.Join(worker.WorkDir, worker.SubmitID)
	)

	lang, err := lang.NewLang(worker.Type, sourcePath, binaryPath)
	if err != nil {
		logrus.Error("New lang failed")
		return nil, nil
	}

	if lang.NeedCompile() {
		compiler, err := NewCompiler(worker, lang)
		if err != nil {
			logrus.Error("New compiler failed")
			return nil, err
		}
		cs, err := compiler.Run()
		if err != nil {
			logrus.Error("Compiler run failed")
			return cs, err
		}
	}

	run, err := NewRunner(worker, lang)

	if err != nil {
		logrus.Error("New runner failed")
		return nil, err
	}
	rs, err := run.Run()

	return rs, nil
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
