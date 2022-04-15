package worker

import (
	"github.com/SXUOJ/judge/pkg/rlimit"
	"github.com/SXUOJ/judge/runner"
	"github.com/sirupsen/logrus"
)

func printResult(rt runner.Result) {
	logrus.Debug(
		"\nstatus:", rt.Status,
		"\nstatus: ", rt.Status,
		"\nexitCode: ", rt.ExitCode,
		"\nerror: ", rt.Error,
		"\ntime: ", rt.Time,
		"\nmemory: ", rt.Memory,
		"\nrunTime: ", rt.RunningTime,
		"\nsetUpTime: ", rt.SetUpTime,
	)
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
