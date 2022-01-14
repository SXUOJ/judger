package cgroup

import (
	"os"
)

//Does the file exist
func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}
