package main

import (
	"./dhcp"
	"flag"
	"fmt"
	"os"
)

func cmdSnoop() {
	var iface string

	flag.StringVar(&iface, "i", "", "network `interface` to use")
	flag.Parse()

	if iface == "" {
		usage(os.Args[1])
		os.Exit(1)
	}

	snoop(iface)
}

func snoop(iface string) {

	mac, err := getMAC(iface)
	checkError(err)

	fmt.Printf("Interface: %s [%s]\n", iface, mac)

	// Set up server
	client, err := dhcp.NewClient()
	checkError(err)
	defer client.Close()

	for {
		o, remote, err := client.Receive(-1)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			continue
		}
		ip := remote.IP.String()
		mac := getMACFromIP(ip)
		fmt.Printf("\n<<< Receive packet from %s (%s)\n", ip, getName(ip))
		fmt.Printf("    MAC address: %s (%s)\n", mac, getVendor(mac))
		showPacket(&o)
	}
}
