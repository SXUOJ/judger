package worker

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
	"github.com/SXUOJ/judge/util"
	"github.com/sirupsen/logrus"
)

type Runner struct {
	count         int
	sampleDir     string
	outputDir     string
	realTimeLimit uint64
	r             sandbox.Runner
}

func NewRunner(worker *Worker, lang lang.Lang) (*Runner, error) {
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

	fileInDir, _ := ioutil.ReadDir(sampleDir)
	for _, fi := range fileInDir {
		if fi.IsDir() {
			continue
		} else {
			sampleCount++
		}
	}

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
	}, nil
}

func (runn *Runner) Run() (RunResults, error) {
	var (
		wg      sync.WaitGroup
		lock    sync.Mutex
		results = make([]RunResult, runn.count)
	)

	for i := 0; i < runn.count; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			logrus.Debug("Runner run start")

			sampleIdStr := strconv.FormatInt(int64(id), 10)

			input := strings.Join([]string{sampleIdStr, "in"}, ".")
			inputFileName := filepath.Join(runn.sampleDir, input)

			output := strings.Join([]string{sampleIdStr, "out"}, ".")
			erroR := strings.Join([]string{sampleIdStr, "err"}, ".")
			outputFileName := filepath.Join(runn.outputDir, output)
			errorFileName := filepath.Join(runn.outputDir, erroR)
			file, err := prepareFiles(inputFileName, outputFileName, errorFileName)
			if err != nil {
				logrus.Error("failed to prepare files: %v", err)
				return
			}
			defer closeFiles(file)

			// if not defined, then use the original value
			fds := make([]uintptr, len(file))
			for i, f := range file {
				if f != nil {
					fds[i] = f.Fd()
				} else {
					fds[i] = uintptr(i)
				}
			}

			runResult := RunResult{}

			r := runn.r
			r.Files = fds
			res, err := run(&r, runn.realTimeLimit)
			runResult.SampleId = id
			runResult.Status = convertStatus(res.Status)
			runResult.ExitCode = res.ExitCode
			runResult.Error = res.Error
			runResult.SetUpTime = res.SetUpTime
			runResult.RunningTime = res.RunningTime / time.Millisecond
			runResult.Time = res.Time / time.Millisecond
			runResult.Memory = res.Memory >> 20

			if res.Status != runner.StatusNormal || err != nil {
				logrus.Error("runByOne failed")
				return
			}

			if ok := runn.Compare(sampleIdStr); ok {
				runResult.Status = StatusAC
			} else {
				runResult.Status = StatusWA
			}

			lock.Lock()
			results = append(results, runResult)
			lock.Unlock()
		}(i + 1)
	}
	wg.Wait()
	return results, nil
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
