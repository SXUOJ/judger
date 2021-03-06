package judger

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/SXUOJ/judge/pkg/rlimit"
	"github.com/SXUOJ/judge/runner"
	"github.com/sirupsen/logrus"
)

func writeFile(path string, text []byte) error {
	err := os.WriteFile(path, text, filePerm)
	return err
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func getFileNum(dir string) int {
	var (
		count = 0
	)
	fileInDir, _ := ioutil.ReadDir(dir)
	for _, fi := range fileInDir {
		if fi.IsDir() {
			continue
		} else {
			count++
		}
	}
	return count
}

func parseLimit(timeLimit, realTimeLimit, outputLimit, stackLimit, memoryLimit uint64) (rlimit.RLimits, runner.Limit) {
	rlimits := rlimit.RLimits{
		CPU:         timeLimit,
		CPUHard:     realTimeLimit,
		FileSize:    outputLimit << 20,
		Stack:       stackLimit << 20,
		Data:        memoryLimit << 20,
		OpenFile:    256,
		DisableCore: true,
	}

	limit := runner.Limit{
		TimeLimit:   time.Duration(timeLimit) * time.Second,
		MemoryLimit: runner.Size(memoryLimit << 20),
	}

	return rlimits, limit
}

func printLimit(rl *rlimit.RLimits) {
	logrus.Debug(
		"\ncpu: ", rl.CPU,
		"\ncpuHard: ", rl.CPUHard,
		"\nfileSize: ", rl.FileSize,
		"\nstack: ", rl.Stack,
		"\ndata: ", rl.Data,
		"\nopenfile: ", rl.OpenFile,
		"\ndisableCore", rl.DisableCore,
	)
}
