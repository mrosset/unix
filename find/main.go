package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	debug = flag.Bool("d", false, "output debugging details")
)

func main() {
	flag.Parse()
	args := flag.Args()
	if flag.NArg() == 0 {
		args = []string{"."}
	}
	if err := filepath.Walk(args[0], walkFunc); err != nil {
		log.Fatal(err)
	}
}

func walkFunc(path string, info os.FileInfo, err error) error {
	fmt.Fprintln(os.Stdout, path)
	return nil
}
