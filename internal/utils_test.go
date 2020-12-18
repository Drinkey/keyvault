package internal

import (
	"os"
	"testing"
)

func TestFileExistDirReturnsFalse(t *testing.T) {
	f := "/tmp/somedir"
	os.Mkdir(f, 0777)
	r := FileExist(f)
	if r {
		t.Log("expect return false but got true")
		t.Fail()
	}
	os.Remove(f)
}

func TestFileExistFileReturnsTrue(t *testing.T) {
	f := "/tmp/somefile"
	// file := os.File{}
	os.Create(f)
	r := FileExist(f)
	if !r {
		t.Log("expect return false but got true")
		t.Fail()
	}
	os.Remove(f)
}
