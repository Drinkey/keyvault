package v1

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestPingResponse(t *testing.T) {
	r := gin.Default()
	w := httptest.NewRecorder()
	r.GET("/api/v1/ping", Ping)
	req, _ := http.NewRequest("GET", "/api/v1/ping", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fail()
	}
	if !strings.Contains(w.Body.String(), "pong") {
		t.Fail()
	}
}
