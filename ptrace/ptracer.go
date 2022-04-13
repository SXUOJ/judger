package ptrace

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/Sxu-Online-Judge/judge/runner"
	"golang.org/x/sys/unix"
)

type Ptracer struct {
	Handler
	Runner
	runner.Limit
}

type Handler interface {
	Handle(*Context) TraceAction
	Debug(v ...interface{})
}

type Runner interface {
	Start() (int, error)
}

func (p *Ptracer) Trace(c context.Context) (result runner.Result) {
	// ptrace is thread based (kernel proc)
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// Start the runner
	pgid, err := p.Runner.Start()
	p.Handler.Debug("tracer started: ", pgid, err)
	if err != nil {
		p.Handler.Debug("start tracee failed: ", err)
		result.Status = runner.StatusSystemError
		result.Error = err.Error()
		return
	}
	return p.trace(c, pgid)
}

func (p *Ptracer) trace(c context.Context, pgid int) (result runner.Result) {
	cc, cancel := context.WithCancel(c)
	defer cancel()

	// handle cancelation
	go func() {
		<-cc.Done()
		killAll(pgid)
	}()

	sTime := time.Now()
	ph := newPtraceHandle(p, pgid)

	// handler potential panic and tle
	// also ensure processes was well terminated
	defer func() {
		if err := recover(); err != nil {
			p.Handler.Debug("panic: ", err)
			result.Status = runner.StatusSystemError
			result.Error = fmt.Sprintf("%v", err)
		}
		// kill all tracee upon return
		killAll(pgid)
		collectZombieProcess(pgid)
		if !ph.fTime.IsZero() {
			result.SetUpTime = ph.fTime.Sub(sTime)
			result.RunningTime = time.Since(ph.fTime)
		}
	}()

	// ptrace pool loop
	for {
		var (
			wstatus unix.WaitStatus // wait4 wait status
			rusage  unix.Rusage     // wait4 rusage
			pid     int             // store pid of wait4 result
			err     error
		)
		if ph.execved {
			// Wait for all child in the process group
			pid, err = unix.Wait4(-pgid, &wstatus, unix.WALL, &rusage)
		} else {
			// Ensure the process have called setpgid
			pid, err = unix.Wait4(pgid, &wstatus, unix.WALL, &rusage)
		}
		if err == unix.EINTR {
			p.Handler.Debug("wait4 EINTR")
			continue
		}
		if err != nil {
			p.Handler.Debug("wait4 failed: ", err)
			result.Status = runner.StatusSystemError
			result.Error = err.Error()
			return
		}
		p.Handler.Debug("---pid: ", pid)

		// update rusage
		if pid == pgid {
			userTime, userMem, curStatus := p.checkUsage(rusage)
			result.Status = curStatus
			result.Time = userTime
			result.Memory = userMem
			if curStatus != runner.StatusNormal {
				return
			}
		}

		status, exitStatus, errStr, finished := ph.handle(pid, wstatus)
		if finished || status != runner.StatusNormal {
			result.Status = status
			result.ExitCode = exitStatus
			result.Error = errStr
			return
		}
	}
}

func (p *Ptracer) checkUsage(rusage unix.Rusage) (time.Duration, runner.Size, runner.Status) {
	status := runner.StatusNormal
	// update resource usage and check against limits
	userTime := time.Duration(rusage.Utime.Nano()) // ns
	userMem := runner.Size(rusage.Maxrss << 10)    // bytes

	// check tle / mle
	if userTime > p.Limit.TimeLimit {
		status = runner.StatusTimeLimitExceeded
	}
	if userMem > p.Limit.MemoryLimit {
		status = runner.StatusMemoryLimitExceeded
	}
	return userTime, userMem, status
}

func killAll(pgid int) {
	unix.Kill(-pgid, unix.SIGKILL)
}

func collectZombieProcess(pgid int) {
	var wstatus unix.WaitStatus
	for {
		if _, err := unix.Wait4(-pgid, &wstatus, unix.WALL|unix.WNOHANG, nil); err != unix.EINTR && err != nil {
			break
		}
	}
}
