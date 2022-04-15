package web

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SXUOJ/judge/model"
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

	params := model.Submit{}

	paramsByte, _ := json.Marshal(params)
	// log.Printf("%s\n", paramsByte)
	req, _ := http.NewRequest("POST", "/submit", bytes.NewBuffer(paramsByte))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	log.Println(w.Body.String())
	assert.Equal(t, 200, w.Code)
}
