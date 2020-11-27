package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Drinkey/keyvault/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPingResponse(t *testing.T) {
	Convey("/ping function", t, func() {
		r := setupRouter()
		w := httptest.NewRecorder()
		Convey("GET /ping should return pong in response body", func() {
			req, _ := http.NewRequest("GET", "/v1/ping", nil)
			r.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, http.StatusOK)
			So(w.Body.String(), ShouldContainSubstring, "pong")
		})
	})
}

func TestNamespace(t *testing.T) {
	Convey("/vault function", t, func() {
		uri := "/v1/vault"
		r := setupRouter()
		w := httptest.NewRecorder()
		var test_ns_id int

		Convey("POST /vault: create new namespace should success", func() {
			reqJson, _ := json.Marshal(map[string]string{
				"name": "TEST_NS",
			})
			req, _ := http.NewRequest("POST", uri, bytes.NewBuffer(reqJson))
			r.ServeHTTP(w, req)
			var respNamespace model.Namespace
			json.Unmarshal(w.Body.Bytes(), &respNamespace)
			test_ns_id = respNamespace.ID

			So(w.Code, ShouldEqual, http.StatusCreated)
			So(respNamespace.ID, ShouldEqual, "TEST_NS")

		})
		Convey("POST /vault: create namespace with duplicated name should fail", func() {
			reqJson, _ := json.Marshal(map[string]string{
				"name": "TEST_NS",
			})
			req, _ := http.NewRequest("POST", uri, bytes.NewBuffer(reqJson))
			r.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, http.StatusBadRequest)
			So(w.Body.String(), ShouldContainSubstring, "failed")
		})
		Convey("GET /vault: get all namespace should success", func() {
			req, _ := http.NewRequest("GET", uri, nil)
			r.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, http.StatusOK)
			So(w.Body.String(), ShouldContainSubstring, "TEST_NS")
		})
		Convey("DELETE /vault/:id: delete the namespace just created should success", func() {
			req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%d", uri, test_ns_id), nil)
			r.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, http.StatusOK)
			So(w.Body.String(), ShouldContainSubstring, "TEST_NS")
		})
	})
}
