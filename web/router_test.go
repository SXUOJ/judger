package web

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SXUOJ/judge/worker"
	"gotest.tools/assert"
)

func TestPingRoute(t *testing.T) {
	router := loadRouter()

	req, _ := http.NewRequest("GET", "/ping", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"msg\":\"pong\"}", w.Body.String())
}

func TestSubmitRouteParameter(t *testing.T) {
	router := loadRouter()

	params := Submit{
		SubmitID:  "001",
		ProblemID: "SXU001",
		FileName:  "001.c",
		Type:      "C",
		AllowProc: false,

		TimeLimit:     1,
		RealTimeLimit: 1,
		MemoryLimit:   256,
		OutputLimit:   256,
		StackLimit:    256,
	}

	paramsByte, _ := json.Marshal(params)
	// log.Printf("%s\n", paramsByte)
	req, _ := http.NewRequest("POST", "/submit", bytes.NewBuffer(paramsByte))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	log.Println(w.Body.String())
	assert.Equal(t, 200, w.Code)
}

func printResult(rt *worker.RunResult) {
	log.Println(
		"\n---------------",
		"\nsampleId:", rt.SampleId,
		"\nstatus:", rt.Status,
		"\nexitCode: ", rt.ExitCode,
		"\nerror: ", rt.Error,
		"\ntime: ", rt.Time,
		"\nmemory: ", rt.Memory,
		"\nrunTime: ", rt.RunningTime,
		"\nsetUpTime: ", rt.SetUpTime,
	)
}
