package main

import (
	"flag"
	"fmt"
	"github.com/str1ngs/util/console"
	"log"
	"os"
	"os/user"
	"syscall"
)

var (
	timeFmt = "Jan _2 15:04"
	long    = flag.Bool("l", false, "use long listing format")
	all     = flag.Bool("a", false, "do not ignore entries staring with .")
)

func main() {
	flag.Parse()
	args := flag.Args()
	if flag.NArg() == 0 {
		args = []string{"."}
	}
	err := ls(args)
	if err != nil {
		log.Fatal(err)
	}
}

func ls(args []string) (err error) {
	files := []os.FileInfo{}
	for _, a := range args {
		fi, err := os.Stat(a)
		if err != nil {
			return err
		}
		if !fi.IsDir() {
			files = append(files, fi)
			continue
		}
		fd, err := os.Open(a)
		if err != nil {
			return err
		}
		fis, err := fd.Readdir(0)
		if err != nil {
			return err
		}
		fd.Close()
		files = append(files, fis...)
	}
	return list(files)
}

func list(files []os.FileInfo) (err error) {
	for _, f := range files {
		if f.Name()[0] == '.' && !*all {
			continue
		}
		if !*long {
			fmt.Printf("%s ", f.Name())
			continue
		}
		stat := f.Sys().(*syscall.Stat_t)
		user, err := user.LookupId(fmt.Sprintf("%v", stat.Uid))
		if err != nil {
			return err
		}
		if *long {
			console.Println(
				f.Mode(),
				stat.Nlink,
				user.Username,
				user.Gid,
				f.Size(),
				f.ModTime().Format(timeFmt),
				f.Name(),
			)
		}
	}
	console.Flush()
	if !*long {
		fmt.Println()
	}
	return
}
