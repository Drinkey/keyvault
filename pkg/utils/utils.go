package utils

import (
	"os"
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
