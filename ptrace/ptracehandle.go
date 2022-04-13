package ptrace

import (
	"time"

	"github.com/Sxu-Online-Judge/judge/runner"
	"golang.org/x/sys/unix"
)

type ptraceHandle struct {
	*Ptracer
	pgid    int
	traced  map[int]bool
	execved bool
	fTime   time.Time
}

func newPtraceHandle(p *Ptracer, pgid int) *ptraceHandle {
	return &ptraceHandle{p, pgid, make(map[int]bool), false, time.Time{}}
}

func (ph *ptraceHandle) handle(pid int, wstatus unix.WaitStatus) (status runner.Status, exitStatus int, errStr string, finished bool) {
	status = runner.StatusNormal
	// check process status
	switch {
	case wstatus.Exited():
		delete(ph.traced, pid)
		ph.Handler.Debug("process exited: ", pid, wstatus.ExitStatus())
		if pid == ph.pgid {
			finished = true
			if ph.execved {
				exitStatus = wstatus.ExitStatus()
				if exitStatus == 0 {
					status = runner.StatusNormal
				} else {
					status = runner.StatusNonzeroExitStatus
				}
				return
			}
			status = runner.StatusSystemError
			errStr = "child process exit before execve"
			return
		}

	case wstatus.Signaled():
		sig := wstatus.Signal()
		ph.Handler.Debug("ptrace signaled: ", sig)
		if pid == ph.pgid {
			delete(ph.traced, pid)
			switch sig {
			case unix.SIGXCPU, unix.SIGKILL:
				status = runner.StatusTimeLimitExceeded
			case unix.SIGXFSZ:
				status = runner.StatusOutputLimitExceeded
			case unix.SIGSYS:
				status = runner.StatusDisallowedSyscall
			default:
				status = runner.StatusSignalled
			}
			exitStatus = int(sig)
			return
		}
		unix.PtraceCont(pid, int(sig))

	case wstatus.Stopped():
		// Set option if the process is newly forked
		if !ph.traced[pid] {
			ph.Handler.Debug("set ptrace option for", pid)
			ph.traced[pid] = true
			// Ptrace set option valid if the tracee is stopped
			if err := setPtraceOption(pid); err != nil {
				status = runner.StatusSystemError
				errStr = err.Error()
				return
			}
		}

		stopSig := wstatus.StopSignal()
		// Check stop signal, if trap then check seccomp
		switch stopSig {
		case unix.SIGTRAP:
			switch trapCause := wstatus.TrapCause(); trapCause {
			case unix.PTRACE_EVENT_SECCOMP:
				if ph.execved {
					// give the customized handle for syscall
					err := ph.handleTrap(pid)
					if err != nil {
						status = runner.StatusDisallowedSyscall
						errStr = err.Error()
						return
					}
				} else {
					ph.Handler.Debug("ptrace seccomp before execve (should be the execve syscall)")
				}

			case unix.PTRACE_EVENT_CLONE:
				ph.Handler.Debug("ptrace stop clone")
			case unix.PTRACE_EVENT_VFORK:
				ph.Handler.Debug("ptrace stop vfork")
			case unix.PTRACE_EVENT_FORK:
				ph.Handler.Debug("ptrace stop fork")
			case unix.PTRACE_EVENT_EXEC:
				// forked tracee have successfully called execve
				if !ph.execved {
					ph.fTime = time.Now()
					ph.execved = true
				}
				ph.Handler.Debug("ptrace stop exec")

			default:
				ph.Handler.Debug("ptrace unexpected trap cause: ", trapCause)
			}
			unix.PtraceCont(pid, 0)
			return

		// check if cpu rlimit hit
		case unix.SIGXCPU:
			status = runner.StatusTimeLimitExceeded
		case unix.SIGXFSZ:
			status = runner.StatusOutputLimitExceeded
		}
		if status != runner.StatusNormal {
			return
		}
		// Likely encountered SIGSEGV (segment violation)
		// Or compiler child exited
		if stopSig != unix.SIGSTOP {
			ph.Handler.Debug("ptrace unexpected stop signal: ", stopSig)
		}
		ph.Handler.Debug("ptrace stopped")
		unix.PtraceCont(pid, int(stopSig))
	}
	return
}

func (ph *ptraceHandle) handleTrap(pid int) error {
	ph.Handler.Debug("seccomp traced")
	if ph.Handler != nil {
		ctx, err := getTrapContext(pid)
		if err != nil {
			return err
		}
		act := ph.Handler.Handle(ctx)

		switch act {
		case ActionTrace:
			// https://www.kernel.org/doc/Documentation/prctl/pkg/seccomp_filter.txt
			return ctx.skipSyscall()

		case ActionKill:
			return runner.StatusDisallowedSyscall
		}
	}
	return nil
}

func setPtraceOption(pid int) error {
	const ptraceFlags = unix.PTRACE_O_TRACESECCOMP | unix.PTRACE_O_EXITKILL | unix.PTRACE_O_TRACEFORK |
		unix.PTRACE_O_TRACECLONE | unix.PTRACE_O_TRACEEXEC | unix.PTRACE_O_TRACEVFORK
	return unix.PtraceSetOptions(pid, ptraceFlags)
}
