package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"util"
	"util/console"
	"util/file"
)

var checkf = util.CheckFatal

func main() {
	flag.Parse()
	args := flag.Args()
	dest := args[len(args)-1]
	args = args[0 : len(args)-1]
	if !file.Exists(dest) {
		err := fmt.Errorf("Directory %s does not exist.", dest)
		checkf(err)
	}
	//fmt.Printf("args: %s -> dest: %s", args, dest)
	for _, arg := range args {
		err := copy(arg, dest)
		checkf(err)
	}
}

func copy(org string, dest string) (err error) {
	st, err := os.Stat(org)
	if err != nil {
		return err
	}
	od, err := os.Open(org)
	if err != nil {
		return err
	}
	defer od.Close()
	dpath := path.Join(dest, path.Base(org))
	dd, err := os.Create(dpath)
	if err != nil {
		return err
	}
	defer dd.Close()
	msg := fmt.Sprintf("cp: %s", path.Base(org))
	pb := console.NewProgressBarWriter(msg, st.Size(), dd)
	_, err = io.Copy(pb, od)
	fmt.Println()
	if err != nil {
		return err
	}
	return nil
}
