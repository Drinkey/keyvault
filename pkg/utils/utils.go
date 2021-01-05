package utils

import (
	"os"
	"strings"
	"time"
)

func FileExist(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func TimeRange(validation int) (time.Time, time.Time) {
	start := time.Now()
	return start, start.AddDate(validation, 0, 0)
}

func DirUpLevel(p string, level int) string {
	if level > 0 {
		panic("level should <= 0")
	}
	dirs := strings.Split(p, "/")
	return strings.Join(dirs[:len(dirs)+level], "/")
}
