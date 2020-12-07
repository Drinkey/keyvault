package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Drinkey/keyvault/controller/namespace"
	"github.com/Drinkey/keyvault/controller/secret"
)

func TestCreateSecretShouldSuccess(t *testing.T) {
	r, w := setup()
	defer teardown()
	uri := "/v1/vault"

	// 1. Create namespace for secret
	r.POST(uri, namespace.Create)
	reqJson, _ := json.Marshal(map[string]string{
		"name": "TEST_NS",
	})
	req, _ := http.NewRequest("POST", uri, bytes.NewBuffer(reqJson))
	r.ServeHTTP(w, req)
	// 2. Create a new secret under the namespace
	r.POST(fmt.Sprintf("%s/:namespace", uri), secret.Create)
	reqJson, _ = json.Marshal(map[string]string{
		"key":   "TEST_SECRET",
		"value": "The_hiddenPassw0rd!",
	})
	req, _ = http.NewRequest("POST", fmt.Sprintf("%s/TEST_NS", uri), bytes.NewBuffer(reqJson))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Logf("actual status code: %d", w.Code)
		t.Fail()
	}
	if !strings.Contains(w.Body.String(), "success") {
		t.Logf("actual body:\n%s\n", w.Body.String())
		t.Fail()
	}
}

func TestGetSecretWithSpecificNamespaceAndKey(t *testing.T) {
	r, w := setup()
	defer teardown()
	uri := "/v1/vault"

	// 1. Create namespace for secret
	r.POST(uri, namespace.Create)
	reqJson, _ := json.Marshal(map[string]string{
		"name": "TEST_NS",
	})
	req, _ := http.NewRequest("POST", uri, bytes.NewBuffer(reqJson))
	r.ServeHTTP(w, req)

	// 2. Create a new secret under the namespace
	r.POST(fmt.Sprintf("%s/:namespace", uri), secret.Create)

	reqJson, _ = json.Marshal(map[string]string{
		"key":   "TEST_SECRET",
		"value": "The_hiddenPassw0rd!",
	})
	req, _ = http.NewRequest("POST", fmt.Sprintf("%s/TEST_NS", uri), bytes.NewBuffer(reqJson))
	req.Header.Add("content-type", `application/json`)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 3. Query the secret with specific key
	r.GET(fmt.Sprintf("%s/:namespace", uri), secret.Query)

	req, _ = http.NewRequest("GET", fmt.Sprintf("%s/TEST_NS?q=TEST_SECRET", uri), nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Logf("actual status code: %d", w.Code)
		t.Fail()
	}
	if !strings.Contains(w.Body.String(), "TEST_NS") {
		t.Logf("actual body:\n%s\n", w.Body.String())
		t.Fail()
	}
}
