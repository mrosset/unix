package unix

import (
	"testing"
)

var (
	existsFiles   = []string{"unix.go", "unix_test.go", "find", "goz"}
	notExistFiles = []string{"aaaaaaaa", "bbbbbbbbb"}
)

func TestFileExists(t *testing.T) {
	for _, f := range existsFiles {
		exists := fileExists(f)
		if !exists {
			t.Errorf("expect to find %s got %v", f, exists)
		}
		t.Logf("%s -> %v", f, exists)
	}

	for _, f := range notExistFiles {
		exists := fileExists(f)
		if exists {
			t.Errorf("expect not to find %s got %v", f, exists)
		}
		t.Logf("%s -> %v", f, exists)
	}
}
