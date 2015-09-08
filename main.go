package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

type command struct {
	name   string
	handle func()
}

var cmd []command

func init() {
	cmd = []command{
		{"discover", cmdDiscover},
	}
}

func cmdDiscover() {
	var iface string
	var secs int

	flag.StringVar(&iface, "i", "", "network `interface` to use")
	flag.IntVar(&secs, "t", 5, "timeout in seconds")
	flag.Parse()

	if iface == "" {
		usage(os.Args[1])
		os.Exit(1)
	}

	timeout := time.Duration(secs) * time.Second

	discover(iface, timeout)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}

func usage(c string) {

	cc := c
	if c == "" {
		cc = "<command>"
	}

	fmt.Fprintf(os.Stderr, "usage: %s %s [options]\n", os.Args[0], cc)

	if c == "" {
		fmt.Fprintf(os.Stderr, "available commands:")
		for _, c := range cmd {
			fmt.Fprintf(os.Stderr, " %s", c.name)
		}
		fmt.Fprintf(os.Stderr, "\n")
	}

	flag.PrintDefaults()
}

func main() {
	if len(os.Args) < 2 {
		usage("")
		os.Exit(1)
	}
	for _, c := range cmd {
		if os.Args[1] == c.name {
			// remove command from argument list
			if len(os.Args) > 2 {
				os.Args = append(os.Args[:1], os.Args[2:]...)
			}
			c.handle()
			os.Exit(0)
		}
	}

	fmt.Fprintf(os.Stderr, "%s: %s: invalid command\n", os.Args[0], os.Args[1])
	os.Exit(1)

}
