package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"syscall"
	"github.com/str1ngs/util/console"
)

var (
	longFmt = "%s\t%s\t%s\t%s\t%s\t%s\t%s\n"
	timeFmt = "Jan _2 15:04"
	long    = flag.Bool("l", false, "use long listing format")
)

func main() {
	flag.Parse()
	path := ""
	switch flag.NArg() {
	case 0:
		path = "."
	default:
		path = flag.Arg(0)
	}
	ls(path)
}

func ls(path string) (err error) {
	fd, err := os.Open(path)
	if err != nil {
		return
	}
	defer fd.Close()
	fis, err := fd.Readdir(0)
	if err != nil {
		return
	}
	for _, f := range fis {
		if !*long {
			fmt.Printf("%s ", f.Name())
			continue
		}
		stat := f.Sys().(*syscall.Stat_t)
		var (
			mode  = fmt.Sprintf("%v", f.Mode())
			nlink = fmt.Sprintf("%v", stat.Nlink)
			uid   = fmt.Sprintf("%v", stat.Uid)
			//gid   = fmt.Sprintf("%v", stat.Gid)
			size  = fmt.Sprintf("%v", f.Size())
			time  = f.ModTime().Format(timeFmt)
			name  = f.Name()
		)
		user, err := user.LookupId(uid)
		if err != nil {
			return err
		}
		if *long {
			console.Println(mode, nlink, user.Username, user.Gid, size, time, name)
		}
	}
	console.Flush()
	if !*long {
		fmt.Println()
	}
	return
}
