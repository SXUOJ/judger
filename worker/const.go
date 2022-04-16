package worker

import "path/filepath"

var (
	RunDir  = filepath.Join("/", "tmp")
	pathEnv = "PATH=/usr/local/bin:/usr/bin:/bin"
)
