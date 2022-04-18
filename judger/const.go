package judger

import "path/filepath"

var (
	RunDir  = filepath.Join("/", "tmp")
	pathEnv = "PATH=/usr/local/bin:/usr/bin:/bin"
)

const (
	filePerm = 0644
	dirPerm  = 0755
)
