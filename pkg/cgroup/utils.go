package cgroup

import (
	"os"
)

//Does the file exist
func fileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

//Write text to file
func writeFile(path string, text []byte) error {
	err := os.WriteFile(path, text, filePerm)
	return err
}

func readFile(path string) ([]byte, error) {
	text, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return text, nil
}

func remove(path string) error {
	if path != "" {
		return os.Remove(path)
	}
	return nil
}
