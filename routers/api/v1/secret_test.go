package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSecretCreateShouldSuccess(t *testing.T) {
	r := gin.Default()
	w := httptest.NewRecorder()
	uri := "/api/v1/vault"

	// 1. Create namespace for secret
	r.POST(uri, CreateNamespace)
	reqJson, _ := json.Marshal(map[string]string{
		"name": "TEST_NS",
	})
	req, _ := http.NewRequest("POST", uri, bytes.NewBuffer(reqJson))
	r.ServeHTTP(w, req)
	// 2. Create a new secret under the namespace
	r.POST(fmt.Sprintf("%s/:namespace", uri), CreateSecret)
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

func TestSecretGetWithSpecificNamespaceAndKey(t *testing.T) {
	r := gin.Default()
	w := httptest.NewRecorder()
	uri := "/api/v1/vault"

	// 1. Create namespace for secret
	r.POST(uri, CreateNamespace)
	reqJson, _ := json.Marshal(map[string]string{
		"name": "TEST_NS_2",
	})
	req, _ := http.NewRequest("POST", uri, bytes.NewBuffer(reqJson))
	r.ServeHTTP(w, req)
	// 2. Create a new secret under the namespace
	r.POST(fmt.Sprintf("%s/:namespace", uri), CreateSecret)
	reqJson, _ = json.Marshal(map[string]string{
		"key":   "TEST_SECRET_2",
		"value": "The_hiddenPassw0rd!",
	})
	req, _ = http.NewRequest("POST", fmt.Sprintf("%s/TEST_NS_2", uri), bytes.NewBuffer(reqJson))
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

	// 3. Query the secret with specific key
	r.GET(fmt.Sprintf("%s/:namespace", uri), GetSecret)

	req, _ = http.NewRequest("GET", fmt.Sprintf("%s/TEST_NS_2?q=TEST_SECRET_2", uri), nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	t.Log("Got Response body:")
	t.Log(w.Body.String())

	if w.Code != http.StatusOK {
		t.Logf("actual status code: %d", w.Code)
		t.Fail()
	}
	if !strings.Contains(w.Body.String(), "The_hiddenPassw0rd!") {
		t.Log("expect to see password but not found")
		t.Logf("actual body:\n%s\n", w.Body.String())
		t.Fail()
	}
	if strings.Contains(w.Body.String(), "TEST_SECRET_2") {
		t.Log("expect not to see secret key but found")
		t.Logf("actual body:\n%s\n", w.Body.String())
		t.Fail()
	}
}
