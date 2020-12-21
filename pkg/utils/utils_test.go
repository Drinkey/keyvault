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
