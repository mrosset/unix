package unix

import (
	"os"
)

func fileExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	if !fi.IsDir() || fi.IsDir() || fi.Mode() == os.ModeSymlink {
		return true
	}
	return false
}
