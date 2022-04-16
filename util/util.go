package util

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"unicode"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func GetFileNum(dir string) int {
	var (
		count = 0
	)
	fileInDir, _ := ioutil.ReadDir(dir)
	for _, fi := range fileInDir {
		if fi.IsDir() {
			continue
		} else {
			count++
		}
	}
	return count
}

func GetFileSize(filename string) int64 {
	var result int64
	filepath.Walk(filename, func(path string, f os.FileInfo, err error) error {
		result = f.Size()
		return nil
	})
	return result
}

func GetMaxInt64(x int64, y int64) int64 {
	if x >= y {
		return x
	} else {
		return y
	}
}

func SplitCmd(s string) (res []string) {
	var buf bytes.Buffer
	insideQuotes := false
	for _, r := range s {
		switch {
		case unicode.IsSpace(r) && !insideQuotes:
			if buf.Len() > 0 {
				res = append(res, buf.String())
				buf.Reset()
			}
		case r == '"' || r == '\'':
			if insideQuotes {
				res = append(res, buf.String())
				buf.Reset()
				insideQuotes = false
				continue
			}
			insideQuotes = true
		default:
			buf.WriteRune(r)
		}
	}
	if buf.Len() > 0 {
		res = append(res, buf.String())
	}
	return
}
