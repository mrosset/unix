package main

import (
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	*long = true
	os.Args = append(os.Args, "/home/strings/", ".")
	main()
}
