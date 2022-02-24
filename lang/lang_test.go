package lang

import (
	"fmt"
	"testing"
)

func TestEveryType(t *testing.T) {
	langs := []string{"C", "Cpp", "Go", "Python2", "Python3"}

	for i, v := range langs {
		lang, err := NewLang(v, "dir")
		must(err, t)

		fmt.Printf("%d: Need Compile: %v, Compile: %v, Run: %v\n", i+1, lang.NeedCompile(), lang.Compile(), lang.Run())
	}
}

func must(err error, t *testing.T) {
	if err != nil {
		t.Fail()
	}
}
