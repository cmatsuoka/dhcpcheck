package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [options]\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	var iface string
	var secs int

	flag.StringVar(&iface, "i", "", "network `interface` to use")
	flag.IntVar(&secs, "t", 5, "timeout in seconds")
	flag.Parse()

	if iface == "" {
		usage()
		os.Exit(1)
	}

	timeout := time.Duration(secs) * time.Second

	discover(iface, timeout)
}
