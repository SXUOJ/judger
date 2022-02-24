package test

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/Sxu-Online-Judge/judge/judge"
)

const (
	cTestDir   = "resources/c"
	cppTestDir = "resource/cpp"
)

func TestCompile(t *testing.T) {
	wd, _ := os.Getwd()

	compiler := judge.NewCompile("C", filepath.Join(wd, cTestDir))
	result := compiler.Compile()
	log.Println(result)
}
