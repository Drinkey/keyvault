package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Drinkey/keyvault/controller/namespace"
	"github.com/Drinkey/keyvault/model"
	// . "github.com/smartystreets/goconvey/convey"
)

func TestNamespaceCreationShouldSuccess(t *testing.T) {
	r, w := setup()
	defer teardown()
	uri := "/v1/vault"
	r.POST(uri, namespace.Create)

	reqJson, _ := json.Marshal(map[string]string{
		"name": "TEST_NS",
	})
	req, _ := http.NewRequest("POST", uri, bytes.NewBuffer(reqJson))
	r.ServeHTTP(w, req)

	var respNamespace model.Namespace
	json.Unmarshal(w.Body.Bytes(), &respNamespace)

	if w.Code != http.StatusCreated {
		t.Fail()
	}
	if respNamespace.Name != "TEST_NS" {
		t.Fail()
	}
}

func TestNamespaceCreateDuplicatedShouldFail(t *testing.T) {
	r, w := setup()
	defer teardown()
	uri := "/v1/vault"
	r.POST(uri, namespace.Create)

	reqJson, _ := json.Marshal(map[string]string{
		"name": "TEST_NS",
	})
	req, _ := http.NewRequest("POST", uri, bytes.NewBuffer(reqJson))
	r.ServeHTTP(w, req)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Logf("actual status code: %d", w.Code)
		t.Fail()
	}
	if strings.Contains(w.Body.String(), "failed") {
		t.Logf("actual body:\n%s\n", w.Body.String())
		t.Fail()
	}
}
