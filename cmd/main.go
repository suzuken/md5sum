package main

import (
	"flag"
	"fmt"
	"github.com/suzuken/md5sum"
	"os"
)

var (
	check = flag.Bool("check", false, "check mode")
)

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Println("please specify path to file. usage: md5sum path/to/file")
		return
	}
	p := flag.Arg(0)
	if *check {
		if _, err := md5sum.CheckGlob(p, os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}
	if err := md5sum.ChecksumGlob(p, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
