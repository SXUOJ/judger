package judge

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/Sxu-Online-Judge/judge/lang"
	"github.com/Sxu-Online-Judge/judge/model"
	"github.com/Sxu-Online-Judge/judge/pkg/ptrace"
	"github.com/Sxu-Online-Judge/judge/util"
)

type RunResult struct {
	status     model.JudgeStatus
	timeUsed   int64
	memoryUsed int64

	fileName map[string]int64
}

type Runner struct {
	codeType string
	codeDir  string
	dataDir  string

	limit model.Limit
}

type DoNothing struct{}

func (r *Runner) Run() (runReulst RunResult) {
	files, err := ioutil.ReadDir(r.dataDir)
	if err != nil {
		fmt.Println("Read datadir filed: ", err)
		runReulst.status = model.SystemError
		return
	}

	var timeTotal int64 = 0
	var memoryMax int64 = 0

	inputFileCount := 0
	runReulst.fileName = make(map[string]int64)

	for _, file := range files {
		if strings.Contains(file.Name(), ".in") {
			temp := make(chan RunResult)

			go func() {
				tp := runByOneFile(r.codeType, r.codeDir, r.dataDir, file, r.limit)
				temp <- tp
			}()

			tempRr := <-temp
			inputFileCount++

			runReulst.fileName[file.Name()] = int64(tempRr.status)
			if tempRr.status != model.StatusAC {
				runReulst.status = tempRr.status
				runReulst.timeUsed = tempRr.timeUsed
				runReulst.memoryUsed = tempRr.memoryUsed
				return
			}

			timeTotal = timeTotal + tempRr.timeUsed
			memoryMax = util.GetMaxInt64(memoryMax, tempRr.memoryUsed)
		}
	}

	if inputFileCount == 0 {
		runReulst.status = model.StatusTLE
		return
	}

	runReulst.status = model.StatusAC
	return RunResult{}
}

func runByOneFile(codeType, codeDir, dataDir string, f os.FileInfo, limit model.Limit) (runResult RunResult) {
	outputFileName := strings.Replace(f.Name(), ".in", ".out", -1)

	inputFile, err := os.Open(filepath.Join(dataDir, f.Name()))
	defer inputFile.Close()
	if err != nil {
		fmt.Println("runByOneFile open inoutFile failed: ", err)
		runResult.status = model.SystemError
		return
	}

	fileName := strings.Replace(f.Name(), ".in", ".user", -1)

	outputFile, err := os.OpenFile(filepath.Join(dataDir, fileName), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
	defer outputFile.Close()
	if err != nil {
		fmt.Println("runByOneFile open outputFile failed: ", err)
		runResult.status = model.SystemError
		return
	}

	lang, err := lang.NewLang(codeType, codeDir)
	if err != nil {
		runResult.status = model.SystemError
		return
	}

	runCmd := lang.Run()
	runCmd.Stdin = inputFile
	runCmd.Stdout = outputFile
	// runCmd.Dir = codeDir
	runCmd.SysProcAttr = &syscall.SysProcAttr{
		Ptrace: true,
	}
	err = runCmd.Start()
	if err != nil {
		fmt.Println("start failed: ", err)
	}

	var timeUsage int64
	var memUsage int64

	runCmd.Wait()

	var regs syscall.PtraceRegs
	pid := runCmd.Process.Pid
	exit := true

	standardOutputFileSize := util.GetFileSize(filepath.Join(dataDir, outputFileName))

	for {
		if exit {
			err = syscall.PtraceGetRegs(pid, &regs)
			isAllowSysCall := ptrace.IsAllowSysCall(regs.Orig_rax)
			if isAllowSysCall != model.Normal && codeType != "Java" {
				runResult.status = isAllowSysCall
				return
			}

			if err != nil {
				break
			}
		}
		err = syscall.PtraceSyscall(pid, 0)
		if err != nil {
			break
		}

		_, err = syscall.Wait4(pid, nil, 0, nil)

		ok, res, runningTime, cpuTime := getResourceUsage(runCmd.Process.Pid)

		if !ok {
			break
		}

		timeUsage = util.GetMaxInt64(timeUsage, cpuTime)
		if timeUsage > limit.TimeLimit || runningTime > 3*limit.TimeLimit {
			runCmd.Process.Kill()
			cmd2 := exec.Command("kill", "-9", strconv.Itoa(pid))
			_ = cmd2.Run()
			runResult.status = model.StatusTLE
			runResult.timeUsed = timeUsage
			runResult.memoryUsed = memUsage
			return
		}

		memUsage = util.GetMaxInt64(memUsage, res)
		if memUsage*3 > limit.MemoryLimit*2 {
			runCmd.Process.Kill()
			cmd2 := exec.Command("kill", "-9", strconv.Itoa(pid))
			_ = cmd2.Run()
			runResult.status = model.StatusMLE
			runResult.timeUsed = timeUsage
			runResult.memoryUsed = memUsage
			return
		}

		outputFileSize := util.GetFileSize(filepath.Join(codeDir, fileName))
		if outputFileSize > 3*standardOutputFileSize {
			runCmd.Process.Kill()
			cmd2 := exec.Command("kill", "-9", strconv.Itoa(pid))
			_ = cmd2.Run()
			runResult.status = model.StatusOLE
			runResult.timeUsed = timeUsage
			runResult.memoryUsed = memUsage
			return
		}

		exit = !exit
	}

	cmd2 := exec.Command("kill", "-9", strconv.Itoa(pid))
	_ = cmd2.Run()

	runResult.timeUsed = timeUsage
	runResult.memoryUsed = memUsage

	if strings.Compare(f.Name(), "000000.in") != 0 {
		//TODO: compare out && std
		runResult.status = model.StatusAC
	} else {
		runResult.status = model.StatusAC
	}

	return
}
