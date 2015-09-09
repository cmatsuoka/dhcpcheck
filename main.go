package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

var cmd map[string]func()

func init() {
	cmd = map[string]func(){
		"discover": cmdDiscover,
		"snoop":    cmdSnoop,
	}
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
		fmt.Fprintf(os.Stderr, "available commands: ")
		var keys []string
		for key, _ := range cmd {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		fmt.Fprintf(os.Stderr, "%s\n", strings.Join(keys, " "))
	}

	flag.PrintDefaults()
}

func main() {
	if len(os.Args) < 2 {
		usage("")
		os.Exit(1)
	}

	if handle := cmd[os.Args[1]]; handle != nil {
		// remove command from argument list
		if len(os.Args) > 2 {
			os.Args = append(os.Args[:1], os.Args[2:]...)
		}
		handle()
		os.Exit(0)
	}

	fmt.Fprintf(os.Stderr, "%s: %s: invalid command\n", os.Args[0], os.Args[1])
	os.Exit(1)

}
