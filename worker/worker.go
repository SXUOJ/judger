package worker

import (
	"net/http"

	"github.com/SXUOJ/judge/lang"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Worker struct {
	C    *gin.Context
	Done chan bool
	Lang lang.Lang
	Compiler
	Runner
}

type Compiler interface {
	Start(chan interface{})
}
type Runner interface {
	Start(chan interface{}, chan bool)
}

func (w *Worker) Run() {
	if w.Lang.NeedCompile() {
		resChan := make(chan interface{})
		go w.Compiler.Start(resChan)

		compiling := true
		for compiling {
			select {
			case res := <-resChan:
				if res != nil {
					msg := "Compile Error"
					logrus.Error(msg)
					w.C.JSON(http.StatusOK, gin.H{
						"msg":    msg,
						"result": res,
					})
					w.stop()
					return
				} else {
					logrus.Info("Compile success")
					compiling = false
				}
			}
		}

	}

	done := make(chan bool)
	resChan := make(chan interface{})
	go w.Runner.Start(resChan, done)

	running := true
	for running {
		select {
		case <-done:
			running = false
		case res := <-resChan:
			w.C.JSON(http.StatusOK, gin.H{
				"msg":    "ok",
				"result": res,
			})
		}
	}
	w.stop()
}

func (w *Worker) stop() {
	w.Done <- true
}

func (w *Worker) Wait() {
	for {
		select {
		case <-w.Done:
			return
		}
	}
}
