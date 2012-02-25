package main

import (
	"compress/gzip"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	EXT = ".goz"
)

var (
	fdebug   = flag.Bool("d", false, "output debugging details")
	fverbose = flag.Bool("v", false, "verbose")
	usage    = `Usage: goz command [arguements]
`
	GOPATH   = os.Getenv("GOPATH")
	vlog     = log.New(os.Stdout, "goz: ", 0)
)

func main() {
	flag.Parse()
	cmd, args := poparg(flag.Args())
	switch cmd {
	case "run":
		checkArgs(run, args)
	case "compress":
		check(compress)
	default:
		checkArgs(run, args)
	}
}

func compress() (err error) {
	glob := filepath.Join(GOPATH, "bin", "*") // $GOPATH/bin/*
	info("compressing binaries in", filepath.Dir(glob))
	files, err := filepath.Glob(glob)
	if err != nil {
		return err
	}
	for _, f := range files {
		if f[len(f)-len(EXT):] == EXT {
			verbose("skipping", f)
			continue
		}
		if filepath.Base(f) == "goz" {
			verbose("skipping", f)
			continue
		}
		fd, err := os.Open(f)
		if err != nil {
			return err
		}
		defer fd.Close()
		gzfd, err := os.Create(filepath.Join(GOPATH, "bin", filepath.Base(f)+EXT))
		if err != nil {
			return err
		}
		gz, err := gzip.NewWriterLevel(gzfd, gzip.BestCompression)
		if err != nil {
			return err
		}
		defer gz.Close()
		_, err = io.Copy(gz, fd)
		if err != nil {
			return err
		}
		nfi, err := gzfd.Stat()
		if err != nil {
			return err
		}
		infof("compressed %-8.8s %v", filepath.Base(f), nfi.Size())
	}
	return
}

func run(args []string) (err error) {
	gozpath := filepath.Join(GOPATH, "bin", args[0]+EXT)
	info("running", args)
	dir, err := ioutil.TempDir("", "goz-")
	if err != nil {
		return err
	}
	defer func() {
		info("cleaning", dir)
		err := os.RemoveAll(dir)
		if err != nil {
			vlog.Fatal(err)
		}
	}()
	fd, err := ioutil.TempFile(dir, "")
	if err != nil {
		return err
	}
	defer fd.Close()
	gzfd, err := os.Open(gozpath)
	if err != nil {
		return err
	}
	defer gzfd.Close()
	gz, err := gzip.NewReader(gzfd)
	if err != nil {
		return err
	}
	defer gz.Close()
	_, err = io.Copy(fd, gz)
	if err != nil {
		return err
	}
	if err = fd.Chmod(os.FileMode(0700)); err != nil {
		return err
	}
	fd.Close()
	cmd := exec.Command(fd.Name(), args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func shift(a []string) []string {
	return append(a[:0], a[0+1:]...)
}

func poparg(a []string) (string, []string) {
	if len(a) == 0 {
		flag.Usage()
		log.Fatal()
	}
	if len(a) == 1 {
		return a[0], a
	}
	return a[0], a[1:]
}

func check(fn func() error) {
	if err := fn(); err != nil {
		log.Fatal(err)
	}
}

func checkArgs(fn func([]string) error, args []string) {
	if err := fn(args); err != nil {
		vlog.Fatal(err)
	}
}

func verbose(a ...interface{}) {
	if *fverbose {
		vlog.Println(a...)
	}
}

func info(a ...interface{}) {
	vlog.Println(a...)
}

func infof(f string, a ...interface{}) {
	vlog.Printf(f, a...)
}
