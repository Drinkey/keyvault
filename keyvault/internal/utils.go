package internal

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

// TODO: implement the master key generation
func GenerateMasterKey() string {
	// keyLen := 24
	return "xDeifu-fkeI19-vs313dR"
}
