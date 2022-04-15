package config

import (
	"github.com/SXUOJ/judge/pkg/seccomp"
	"github.com/SXUOJ/judge/sandbox/handle"
)

// SyscallConfig defines extra syscallConfig apply to program type
type SyscallConfig struct {
	DefaultAction seccomp.Action
	Allow         []string
	Trace         []string
}

func GetConf(pType string, allowProc bool) (seccomp.Action, []string, []string, *handle.Handler) {
	var (
		sc            = handle.NewSyscallCounter()
		allow         = append(append([]string{}, defaultSyscallAllows...))
		trace         = append(append([]string{}, defaultSyscallTraces...))
		defaultAction seccomp.Action
	)

	if c, o := runptraceConfig[pType]; o {
		defaultAction = c.DefaultAction
		allow = append(allow, c.Allow...)
		trace = append(trace, c.Trace...)
	}
	if allowProc {
		allow = append(allow, defaultProcSyscalls...)
	}
	allow, trace = cleanTrace(allow, trace)

	return defaultAction, allow, trace, &handle.Handler{
		SyscallCounter: sc,
	}
}

func keySetToSlice(m map[string]bool) []string {
	rt := make([]string, 0, len(m))
	for k := range m {
		rt = append(rt, k)
	}
	return rt
}

func cleanTrace(allow, trace []string) ([]string, []string) {
	// make sure allow, trace no duplicate
	traceMap := make(map[string]bool)
	for _, s := range trace {
		traceMap[s] = true
	}
	allowMap := make(map[string]bool)
	for _, s := range allow {
		if !traceMap[s] {
			allowMap[s] = true
		}
	}
	return keySetToSlice(allowMap), keySetToSlice(traceMap)
}
