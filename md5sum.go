package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Println("please specify path to file. usage: md5sum path/to/file")
		return
	}
	p := flag.Arg(0)
	b, err := ioutil.ReadFile(p)
	if err != nil {
		fmt.Printf("cannot open file: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("%x  %s\n", md5.Sum(b), p)
}
