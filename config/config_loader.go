package config

import "github.com/SXUOJ/judge/sandbox/handle"

// GetConf return file access check set, syscall counter, allow and traced syscall arrays and new args
func GetConf(pType string, allowProc bool) ([]string, []string, *handle.Handler) {
	var (
		sc    = handle.NewSyscallCounter()
		allow = append(append([]string{}, defaultSyscallAllows...))
		trace = append(append([]string{}, defaultSyscallTraces...))
	)

	if c, o := runptraceConfig[pType]; o {
		allow = append(allow, c.Syscall.Allow...)
		trace = append(trace, c.Syscall.Trace...)

	}
	if allowProc {
		allow = append(allow, defaultProcSyscalls...)
	}
	allow, trace = cleanTrace(allow, trace)

	return allow, trace, &handle.Handler{
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
