package judger

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/SXUOJ/judge/config"
	"github.com/SXUOJ/judge/lang"
	"github.com/SXUOJ/judge/pkg/seccomp"
	"github.com/SXUOJ/judge/runner"
	"github.com/SXUOJ/judge/sandbox"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RunResults []RunResult

type RunResult struct {
	SampleId int `json:"sample_id"`
	runner.Result
}

type thisRunner struct {
	submitID      string
	workDir       string
	count         int
	realTimeLimit uint64
	r             sandbox.Runner
	c             *gin.Context
}

func newRunner(submitID, codeType, workDir string, allowProc bool, slimit Limit, lang lang.Lang, c *gin.Context) *thisRunner {
	defaultAction, allow, trace, h := config.GetConf(strings.Join([]string{codeType, "run"}, "-"), allowProc)

	// limit
	rlimits, limit := parseLimit(slimit.TimeLimit, slimit.RealTimeLimit, slimit.OutputLimit, slimit.StackLimit, slimit.MemoryLimit)
	// build seccomp filter
	seccompBuilder := seccomp.Builder{
		Default: defaultAction,
		Allow:   allow,
		Trace:   trace,
	}
	filter, _ := seccompBuilder.Build()

	// load fds
	var (
		sampleDir   = filepath.Join(workDir, "sample")
		sampleCount = 0
	)

	sampleCount = getFileNum(sampleDir)
	return &thisRunner{
		submitID:      submitID,
		workDir:       workDir,
		count:         sampleCount / 2,
		realTimeLimit: slimit.RealTimeLimit,
		r: sandbox.Runner{
			Args: lang.RunArgs(),
			// Env:         []string{pathEnv},
			Env:         os.Environ(),
			WorkDir:     workDir,
			Seccomp:     filter,
			RLimits:     rlimits.PrepareRLimit(),
			Limit:       limit,
			ShowDetails: config.ShowDetails,
			Unsafe:      config.UnSafe,
			Handler:     h,
		},
		c: c,
	}
}

func (thisRun thisRunner) Start(resChan chan interface{}, done chan bool) {
	var (
		wg sync.WaitGroup
	)

	for i := 0; i < thisRun.count; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			var (
				inputFileName  string
				outputFileName string
				errorFileName  string
				files          []*os.File
				sampleDir      = filepath.Join(thisRun.workDir, "sample")
				outputDir      = filepath.Join(thisRun.workDir, "output")
			)

			sampleIdStr := strconv.FormatInt(int64(id), 10)
			input := strings.Join([]string{sampleIdStr, "in"}, ".")
			output := strings.Join([]string{sampleIdStr, "out"}, ".")
			erroR := strings.Join([]string{sampleIdStr, "err"}, ".")
			inputFileName = filepath.Join(sampleDir, input)
			outputFileName = filepath.Join(outputDir, output)
			errorFileName = filepath.Join(outputDir, erroR)
			files, err := prepareFiles(inputFileName, outputFileName, errorFileName)
			if err != nil {
				logrus.Error("failed to prepare files: %v", err)
				return
			}
			defer closeFiles(files)

			// if not defined, then use the original value
			fds := make([]uintptr, len(files))
			for i, f := range files {
				if f != nil {
					fds[i] = f.Fd()
				} else {
					fds[i] = uintptr(i)
				}
			}

			r := thisRun.r
			r.Files = fds
			res, err := run(&r, thisRun.realTimeLimit)
			runResult := convertResult(id, res)

			if res.Status != runner.StatusNormal || err != nil {
				logrus.Error("Program error")
				goto quit
			}

			if ok := thisRun.Compare(sampleIdStr); ok {
				runResult.Status = runner.StatusAccept
			} else {
				runResult.Status = runner.StatusWrongAnswer
			}

		quit:
			resChan <- runResult
			return
		}(i + 1)
	}
	wg.Wait()
	done <- true
}

func (thisRun *thisRunner) Compare(sampleId string) bool {
	//TODO: presentation judge
	var (
		sampleDir = filepath.Join(thisRun.workDir, "sample")
		outputDir = filepath.Join(thisRun.workDir, "output")
	)
	outPath := filepath.Join(outputDir, strings.Join([]string{sampleId, ".out"}, ""))
	ansPath := filepath.Join(sampleDir, strings.Join([]string{sampleId, ".out"}, ""))

	b, err := ioutil.ReadFile(ansPath)
	if err != nil {
		b = []byte{}
	}

	o, err := ioutil.ReadFile(outPath)
	if err != nil {
		o = []byte{}
	}

	ans := plain(b)
	out := plain(o)
	// log.Printf("ans:= %s", ans)
	// log.Printf("out:= %s", out)

	if out == ans {
		return true
	}
	return false
}

func plain(raw []byte) string {
	buf := bufio.NewScanner(bytes.NewReader(raw))
	var b bytes.Buffer
	newline := []byte{'\n'}
	for buf.Scan() {
		b.Write(bytes.TrimSpace(buf.Bytes()))
		b.Write(newline)
	}
	return b.String()
}

func convertResult(id int, res *runner.Result) *RunResult {
	return &RunResult{
		SampleId: id,
		Result: runner.Result{
			Status: res.Status,

			SetUpTime:   res.SetUpTime,
			RunningTime: res.RunningTime / time.Millisecond,
			Time:        res.Time / time.Millisecond,
			Memory:      res.Memory >> 20,

			ExitCode: res.ExitCode,
			Error:    res.Error,
		},
	}
}
