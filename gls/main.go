package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
	"strings"
	"syscall"
	"util/console"
	"util/human"
)

const (
	escape  = "\x1b[00;0;%vm%s\x1b[m"
	timeFmt = "Jan _2 15:04"
)

var (
	// flags
	long = flag.Bool("l", false, "use long listing format")
	all  = flag.Bool("a", false, "do not ignore entries staring with .")

	// color
	ls_colors = os.Getenv("LS_COLORS")
	colors    = map[string]string{}
)

func init() {
	for _, j := range strings.Split(ls_colors, ":") {
		kv := strings.Split(j, "=")
		if len(kv) == 2 {
			colors[kv[0]] = kv[1]
		}
	}
}

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
				//f.Size(),
				human.ByteSize(f.Size()),
				f.ModTime().Format(timeFmt),
				getColor(f),
			)
		}
	}
	console.Flush()
	if !*long {
		fmt.Println()
	}
	return
}

func getColor(fi os.FileInfo) string {
	key := fmt.Sprintf("*%s", path.Ext(fi.Name()))
	switch {
	case fi.Mode()&os.ModeDir != 0:
		key = "di"
	case fi.Mode()&os.ModeSymlink != 0:
		key = "ln"
	case colors[key] == "":
		return fi.Name()
	}
	return fmt.Sprintf(escape, colors[key], fi.Name())
}
