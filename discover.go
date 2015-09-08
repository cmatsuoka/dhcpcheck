package main

import (
	"./dhcp"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

func cmdDiscover() {
	var iface string
	var secs int
	var sendOnly bool

	flag.StringVar(&iface, "i", "", "network `interface` to use")
	flag.IntVar(&secs, "t", 5, "timeout in seconds")
	flag.BoolVar(&sendOnly, "s", false, "send discovery only and ignore offers")
	flag.Parse()

	if iface == "" {
		usage(os.Args[1])
		os.Exit(1)
	}

	timeout := time.Duration(secs) * time.Second
	if sendOnly {
		timeout = 0
	}

	discover(iface, timeout)
}

func discover(iface string, timeout time.Duration) {

	mac, err := getMAC(iface)
	checkError(err)

	fmt.Printf("Interface: %s [%s]\n", iface, mac)

	var conn *net.UDPConn
	if timeout > 0 {
		// Set up server
		addr, err := net.ResolveUDPAddr("udp4", ":68")
		checkError(err)
		conn, err = net.ListenUDP("udp4", addr)
		checkError(err)
		defer conn.Close()
	}

	// Send discover packet
	p := dhcp.NewDiscoverPacket()
	p.ParseMAC(mac)

	fmt.Println("\n>>> Send DHCP discover")
	showPacket(p)
	err = dhcp.Broadcast(p)
	checkError(err)

	if timeout > 0 {
		t := time.Now()
		for time.Since(t) < timeout {
			o, remote, err := dhcp.Receive(conn, timeout)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
				break
			}
			if o.Xid == p.Xid {
				fmt.Printf("\n<<< Receive DHCP offer from %s (%s)\n", remote.IP.String(), getName(remote.IP.String()))
				showPacket(&o)
			}
		}
		fmt.Println("No more offers.")
	}
}
