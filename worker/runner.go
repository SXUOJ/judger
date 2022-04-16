package worker

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"net/http"
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
	"github.com/SXUOJ/judge/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RunResults []RunResult

type RunResult struct {
	SampleId int `json:"sample_id"`
	runner.Result
}

type Runner struct {
	count         int
	sampleDir     string
	outputDir     string
	realTimeLimit uint64
	r             sandbox.Runner
	c             *gin.Context
}

func NewRunner(worker *Worker, lang lang.Lang, c *gin.Context) (*Runner, error) {
	defaultAction, allow, trace, h := config.GetConf(strings.Join([]string{worker.Type, "run"}, "-"), worker.AllowProc)

	// limit
	rlimits, limit := parseLimit(worker.TimeLimit, worker.RealTimeLimit, worker.OutputLimit, worker.StackLimit, worker.MemoryLimit)
	// build seccomp filter
	seccompBuilder := seccomp.Builder{
		Default: defaultAction,
		Allow:   allow,
		Trace:   trace,
	}
	filter, _ := seccompBuilder.Build()

	// load fds
	var (
		sampleDir   = filepath.Join(sampleBaseDir, worker.ProblemID)
		outputDir   = filepath.Join(worker.WorkDir)
		sampleCount = 0
	)

	sampleCount = util.GetFileNum(outputDir)

	if ok, err := util.PathExists(outputDir); err != nil {
		logrus.Warn("Check if path exists failed")
	} else {
		if ok {
			logrus.Println("Output dir exists: ", outputDir)
		} else {
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				return nil, err
			}
		}
	}

	return &Runner{
		count:         sampleCount / 2,
		sampleDir:     sampleDir,
		outputDir:     outputDir,
		realTimeLimit: worker.RealTimeLimit,
		r: sandbox.Runner{
			Args: lang.RunArgs(),
			Env:  os.Environ(),
			// ExecFile: execFile,
			// Files: fds,
			// WorkDir:     worker.WorkDir,
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

func (runn *Runner) Run() {
	var (
		wg      sync.WaitGroup
		lock    sync.Mutex
		results RunResults
	)

	for i := 0; i < runn.count; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			logrus.Debug("Runner run start")

			var (
				inputFileName  string
				outputFileName string
				errorFileName  string
				files          []*os.File
			)

			sampleIdStr := strconv.FormatInt(int64(id), 10)
			input := strings.Join([]string{sampleIdStr, "in"}, ".")
			output := strings.Join([]string{sampleIdStr, "out"}, ".")
			erroR := strings.Join([]string{sampleIdStr, "err"}, ".")
			inputFileName = filepath.Join(runn.sampleDir, input)
			outputFileName = filepath.Join(runn.outputDir, output)
			errorFileName = filepath.Join(runn.outputDir, erroR)
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

			r := runn.r
			r.Files = fds
			res, err := run(&r, runn.realTimeLimit)
			runResult := convertResult(id, res)
			defer func() {
				lock.Lock()
				results = append(results, *runResult)
				lock.Unlock()
			}()

			if res.Status != runner.StatusNormal || err != nil {
				logrus.Error("Program error")
				return
			}

			if ok := runn.Compare(sampleIdStr); ok {
				runResult.Status = runner.StatusAccept
			} else {
				runResult.Status = runner.StatusWrongAnswer
			}

		}(i + 1)
	}
	wg.Wait()

	runn.c.JSON(http.StatusOK, gin.H{
		"msg":    "ok",
		"result": results,
	})
	return
}

func (runn *Runner) Compare(sampleId string) bool {
	//TODO: presentation judge
	outPath := filepath.Join(runn.outputDir, strings.Join([]string{sampleId, ".out"}, ""))
	ansPath := filepath.Join(runn.sampleDir, strings.Join([]string{sampleId, ".out"}, ""))

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
