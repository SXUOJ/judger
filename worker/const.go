package worker

import "path/filepath"

var (
	runDir        = filepath.Join(".", "tmp")
	sampleBaseDir = filepath.Join(".", "sample")
	pathEnv       = "PATH=/usr/local/bin:/usr/bin:/bin"
)
