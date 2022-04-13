package sandbox

import (
	"github.com/Sxu-Online-Judge/judge/runner"
)

type Sandbox struct {
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
