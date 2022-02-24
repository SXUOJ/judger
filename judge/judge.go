package judge

import (
	"fmt"

	"github.com/Sxu-Online-Judge/judge/model"
)

type Judger struct {
	preparer Init
	compiler Compiler
	runner   Runner
}

func NewJudger(submit model.Submit, codeDir, dataDir string) *Judger {
	return &Judger{
		preparer: &Preparer{
			codeType:   submit.CodeType,
			sourceCode: submit.CodeSource,
		},
		compiler: Compiler{
			codeType: submit.CodeType,
			codeDir:  codeDir,
		},
		runner: Runner{
			codeType: submit.CodeType,
			codeDir:  codeDir,
			dataDir:  dataDir,

			limit: submit.Limit,
		},
	}
}

func (judger *Judger) Run(limit model.Limit) (result model.Result) {
	initResult := judger.preparer.Init()
	if initResult != model.Normal {
		result.Status = initResult
		return
	}
	fmt.Println("init OK")

	compileResult := judger.compiler.Compile()
	if compileResult.status == model.SystemError {
		result.Status = compileResult.status
		return
	} else if compileResult.status == model.StatusCE {
		result.Status = compileResult.status
		result.ErrorInf = compileResult.compileErrorInfo
		return
	}
	fmt.Println("compile OK")

	runResult := judger.runner.Run()
	result.FileName = make(map[string]int64)

	result.TimeUsed = runResult.timeUsed
	result.MemoryUsed = runResult.memoryUsed
	result.FileName = runResult.fileName

	if runResult.status != model.StatusAC && runResult.status < model.Score {
		result.Status = runResult.status
		return
	}

	return
}
