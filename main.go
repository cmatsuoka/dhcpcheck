package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strings"
)

var (
	stats Statistics
	cmd   map[string]func()
)

type Statistics struct {
	pkrec, pkproc, pksent uint64
	msg, smac, rmac       map[string]uint64
}

func init() {
	stats = Statistics{}
	stats.smac = map[string]uint64{}
	stats.rmac = map[string]uint64{}
	stats.msg = map[string]uint64{}

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

func setupSummary() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		summary()
		os.Exit(1)
	}()
}

func summary() {
	fmt.Println("\nPacket summary")
	fmt.Println("  Packets sent      :", stats.pksent)
	fmt.Println("  Packets received  :", stats.pkrec)
	fmt.Println("  Packets processed :", stats.pkproc)

	fmt.Println("\nMessage Types")
	for key, val := range stats.msg {
		fmt.Printf("  %-12.12s : %d\n", key, val)
	}

	fmt.Println("\nVendor stats")

	vendor := map[string]bool{}

	vsent := map[string]uint64{}
	for key, _ := range stats.smac {
		v := VendorFromMAC(key)
		vsent[v]++
		vendor[v] = true
	}

	vrec := map[string]uint64{}
	for key, _ := range stats.rmac {
		v := VendorFromMAC(key)
		vrec[v]++
		vendor[v] = true
	}

	for key, _ := range vendor {
		fmt.Printf("  %-8.8s : %d out / %d in\n",
			key, vsent[key], vrec[key])
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
		summary()
		os.Exit(0)
	}

	fmt.Fprintf(os.Stderr, "%s: %s: invalid command\n", os.Args[0], os.Args[1])
	os.Exit(1)

}
