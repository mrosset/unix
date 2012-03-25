package main

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	os.Args = append(os.Args, ".")
	fmt.Println("normal")
	main()
	*long = true
	fmt.Println("long")
	main()
}
