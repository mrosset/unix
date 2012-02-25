package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	debug = flag.Bool("d", false, "output debugging details")
	usage = `Usage: find path [arguements]
`
)

var Usage = func() {
	fmt.Fprintf(os.Stderr, usage)
	flag.PrintDefaults()
}

func main() {
	start := time.Now()
	defer func() {
		if *debug {
			fmt.Println(os.Args[0], "done in", time.Now().Sub(start))
		}
	}()
	flag.Usage = Usage
	flag.Parse()
	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	path := flag.Arg(0)
	if err := filepath.Walk(path, walkFunc); err != nil {
		log.Fatal(err)
	}
}

func walkFunc(path string, info os.FileInfo, err error) error {
	fmt.Fprintln(os.Stdout, path)
	return nil
}
