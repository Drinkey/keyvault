package utils

import (
	"os"
	"testing"
)

func TestFileExistDirReturnsFalse(t *testing.T) {
	f := "/tmp/somedir"
	os.Mkdir(f, 0777)
	defer os.Remove(f)
	r := FileExist(f)
	if r {
		t.Log("expect return false but got true")
		t.Fail()
	}
}

func TestFileExistFileReturnsTrue(t *testing.T) {
	f := "/tmp/somefile"
	// file := os.File{}
	os.Create(f)
	defer os.Remove(f)
	r := FileExist(f)
	if !r {
		t.Log("expect return false but got true")
		t.Fail()
	}
}

func TestDirUpLevelUsingFullFileName(t *testing.T) {
	p := "/go/src/keyvault/certio/certificate.go"
	upLevel := DirUpLevel(p, -2)
	if upLevel != "/go/src/keyvault" {
		t.Logf("Actual: %s", upLevel)
		t.Fail()
	}
}

func TestDirUpLevelUsingDirNameOneLevel(t *testing.T) {
	p := "/go/src/keyvault/certio"
	upLevel := DirUpLevel(p, -1)
	if upLevel != "/go/src/keyvault" {
		t.Logf("Actual: %s", upLevel)
		t.Fail()
	}
}

func TestDirUpLevelUsingDirNameMultiLevel(t *testing.T) {
	p := "/go/src/keyvault/pkg/utils/utils.go"
	upLevel := DirUpLevel(p, -3)
	if upLevel != "/go/src/keyvault" {
		t.Logf("Actual: %s", upLevel)
		t.Fail()
	}
}

func TestDirUpLevelPositiveLevelWillPanic(t *testing.T) {
	defer func() {
		recover()
	}()
	p := "/go/src/keyvault/pkg/utils/utils.go"
	_ = DirUpLevel(p, 1)
	// only fail when function is not panic
	t.Fail()
}
