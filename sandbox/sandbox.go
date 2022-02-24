package sandbox

import (
	"os"
	"path/filepath"

	"github.com/Sxu-Online-Judge/judge/judge"
	"github.com/Sxu-Online-Judge/judge/model"
)

type Sandbox interface {
	Run() model.Result
}

type StdSandbox struct {
	TimeLimit   int64
	MemoryLimit int64

	TimeUsed   int64
	MemoryUsed int64
}

const (
	baseDir      = "/home/ther/sandbox/"
	baseJudgeDir = "tmp"
	baseDataDir  = "data"
)

func (s *StdSandbox) Run(submit model.Submit) (result model.Result, err error) {
	codeDir := filepath.Join(baseDir, baseJudgeDir, submit.SubmitId)
	dataDir := filepath.Join(baseDir, baseDataDir, submit.ProblemId)

	judger := judge.NewJudger(submit, codeDir, dataDir)
	if judger == nil {
		result.Status = model.SystemError
		return
	}

	result = judger.Run(
		model.Limit{
			TimeLimit:   submit.TimeLimit,
			MemoryLimit: submit.MemoryLimit,
		},
	)
	_ = os.Remove(codeDir)
	return
}
