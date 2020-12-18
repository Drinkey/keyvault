package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNamespaceCreateSuccess(t *testing.T) {
	r := gin.Default()
	w := httptest.NewRecorder()
	uri := "/api/v1/namespace"
	r.POST(uri, CreateNamespace)

	reqJSON, _ := json.Marshal(map[string]string{
		"name": "TEST_NS",
	})
	req, _ := http.NewRequest("POST", uri, bytes.NewBuffer(reqJSON))
	r.ServeHTTP(w, req)

	var respNamespace Namespace
	t.Log(w.Body.String())
	json.Unmarshal(w.Body.Bytes(), &respNamespace)

	if w.Code != http.StatusCreated {
		t.Logf("response code validation failed: code=%d, expected=%d", w.Code, http.StatusCreated)
		t.Fail()
	}
	if respNamespace.Name != "TEST_NS" {
		t.Logf("response body validation failed: %s", respNamespace.Name)
		t.Log(respNamespace)
		t.Fail()
	}
}

func TestNamespaceListSuccess(t *testing.T) {
	r := gin.Default()
	w := httptest.NewRecorder()
	uri := "/api/v1/namespace"
	r.POST(uri, CreateNamespace)
	r.GET(uri, ListNamespaces)

	for _, ns := range []string{"Test_NS-Ls1", "Test_NS-Ls2", "Test_NS-Ls3"} {
		reqJSON, _ := json.Marshal(map[string]string{
			"name": ns,
		})
		req, _ := http.NewRequest("POST", uri, bytes.NewBuffer(reqJSON))
		r.ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			t.Logf("response code validation failed: code=%d, expected=%d", w.Code, http.StatusCreated)
			t.Fail()
		}
	}
	w = httptest.NewRecorder()
	req, _ := http.NewRequest("GET", uri, nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Logf("response code validation failed: code=%d, expected=%d", w.Code, http.StatusOK)
		t.Fail()
	}

	if !strings.Contains(w.Body.String(), "Test_NS-Ls1") {
		t.Logf("actual body:\n%s\n", w.Body.String())
		t.Fail()
	}
}

func TestNamespaceCreateDuplicatedShouldFail(t *testing.T) {
	r := gin.Default()
	w := httptest.NewRecorder()
	uri := "/api/v1/namespace"
	r.POST(uri, CreateNamespace)

	reqJSON, _ := json.Marshal(map[string]string{
		"name": "TEST_NS_DUP",
	})
	req, _ := http.NewRequest("POST", uri, bytes.NewBuffer(reqJSON))
	r.ServeHTTP(w, req)
	var respNamespace Namespace
	t.Log(w.Body.String())
	json.Unmarshal(w.Body.Bytes(), &respNamespace)

	if w.Code != http.StatusCreated {
		t.Logf("first time response code validation failed: code=%d, expected=%d", w.Code, http.StatusCreated)
		t.Fail()
	}

	// send the same request the second time
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Logf("second time response code validation failed: code=%d, expected=%d", w.Code, http.StatusBadRequest)
		t.Fail()
	}
	if strings.Contains(w.Body.String(), "failed") {
		t.Logf("actual body:\n%s\n", w.Body.String())
		t.Fail()
	}

}
