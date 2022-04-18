package service

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/SXUOJ/judge/main/model"
)

func BenchmarkSubmit(b *testing.B) {
	router := loadRouter()
	testCase := newTestCase()
	b.ResetTimer()
	for i, v := range testCase {
		b.Log("Start: ", i)
		req, _ := http.NewRequest("POST", "/submit", bytes.NewBuffer(v))

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		log.Println(w.Body.String())
	}
}

func newTestCase() (b [][]byte) {
	params := []model.Submit{
		{
			SourceCode: `
#include<stdio.h>
int main(){
	int n;
	scanf("%d", &n);
	printf("%d", n);
	return 0;
}`,
			CodeType:  "C",
			Samples:   newSample(10),
			AllowProc: false,

			TimeLimit:     1,
			RealTimeLimit: 1,
			MemoryLimit:   20,
			OutputLimit:   20,
			StackLimit:    20,
		},
	}

	for _, param := range params {
		paramsByte, _ := json.Marshal(param)
		b = append(b, paramsByte)
	}
	return b
}

func newSample(num int) (samples []model.Sample) {
	for i := 0; i < num; i++ {
		a := strconv.FormatInt(int64(i), 10)
		samples = append(samples, model.Sample{
			In:  a,
			Out: a,
		})
	}
	return samples
}
