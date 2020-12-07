package main

import (
	"net/http"
	"strings"
	"testing"

	"github.com/Drinkey/keyvault/controller"
)

func TestPingResponse(t *testing.T) {
	r, w := setup()
	defer teardown()
	r.GET("/v1/ping", controller.Ping)
	req, _ := http.NewRequest("GET", "/v1/ping", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fail()
	}
	if !strings.Contains(w.Body.String(), "pong") {
		t.Fail()
	}
}
