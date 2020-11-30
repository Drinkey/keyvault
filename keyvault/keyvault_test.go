package main

import (
	"io"
	"log"
	"net/http/httptest"
	"os"

	"github.com/gin-gonic/gin"
)

const TESTDB = "/tmp/vault_test.db"
const TESTDB_BASELINE = "/tmp/vault_test_baseline.db"

// Copy the src file to dst. Any existing file will be overwritten and will not
// copy file attributes.
func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func setup() (*gin.Engine, *httptest.ResponseRecorder) {
	log.SetPrefix("TestCase Setup: ")
	log.Printf("backing up db: %s -> %s", TESTDB, TESTDB_BASELINE)
	Copy(TESTDB, TESTDB_BASELINE)
	log.Print("Getting router")
	r := getRouter()
	w := httptest.NewRecorder()
	return r, w
}

func teardown() {
	log.SetPrefix("TestCase TearDown: ")
	log.Printf("restoring database file %s -> %s", TESTDB_BASELINE, TESTDB)
	Copy(TESTDB_BASELINE, TESTDB)
}
