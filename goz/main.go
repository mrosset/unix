package main

import (
	"compress/gzip"
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"unix"
)

const (
	EXT             = ".goz"
	BINFMT_NAME     = "GOZ"
	BINFMT_DIR      = "/proc/sys/fs/binfmt_misc"
	BINFMT_REGISTER = "/proc/sys/fs/binfmt_misc/register"
	BINFMT_REGFILE  = "/proc/sys/fs/binfmt_misc/GOZ"
	BINFMT_REGFMT   = `:GOZ:M::\x1f\x8b::/home/strings/gocode/bin/goz:`
)

var (
	// flags
	fdebug   = flag.Bool("d", false, "output debugging details")
	fverbose = flag.Bool("v", false, "verbose")
	usage    = `Usage: goz command [arguements]
`
	fpath    = flag.String("path", os.Getenv("GOPATH")+"/bin", "path to compress binaries")

	vlog = log.New(os.Stdout, "goz: ", 0)

	// errors
	NotRootRegisterError = errors.New("you must be root to register with binfmt")
)

func main() {
	flag.Parse()
	cmd := flag.Arg(0)
	switch cmd {
	case "run":
		check(run)
	case "compress":
		check(compress)
	case "register":
		check(register)
	case "unregister":
		check(unregister)
	default:
		err := interp(cmd)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func register() (err error) {
	if os.Getuid() != 0 {
		return NotRootRegisterError
	}
	if isBinFmtRegistered() {
		info(BINFMT_NAME, "is already registered ")
		return nil
	}
	fd, err := os.OpenFile(BINFMT_REGISTER, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer fd.Close()
	if _, err = fd.WriteString(BINFMT_REGFMT); err != nil {
		return err
	}
	info("Registered", BINFMT_NAME, "with binfmt")
	return nil
}

func unregister() (err error) {
	if os.Getuid() != 0 {
		return NotRootRegisterError
	}
	fd, err := os.OpenFile(BINFMT_REGFILE, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer fd.Close()
	if _, err = fd.WriteString("-1"); err != nil {
		return err
	}
	info("Unregisterd", BINFMT_NAME, "with binfmt")
	return nil
}

func isBinFmtRegistered() bool {
	return unix.FileExists(BINFMT_REGFILE)
}

func compress() (err error) {
	glob := filepath.Join(*fpath, "*") // $GOPATH/bin/*
	info("Compressing binaries in", filepath.Dir(glob))
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
		gzfd, err := os.Create(filepath.Join(*fpath, filepath.Base(f)+EXT))
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
		gzfd.Chmod(os.FileMode(0755))
		infof("Compressed %-8.8s %v", filepath.Base(f), nfi.Size())
	}
	return
}

func run() (err error) {
	info("would run here")
	return nil
}

func interp(path string) (err error) {
	info("Running", path, flag.Args()[1:])
	dir, err := ioutil.TempDir("", "goz-")
	if err != nil {
		return err
	}
	defer func() {
		info("Cleaning", dir)
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
	gzfd, err := os.Open(path)
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
	cmd := exec.Command(fd.Name(), flag.Args()[1:]...)
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
