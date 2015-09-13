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

	setupSummary()

	snoop(iface)
}

type message struct {
	origin string
	packet dhcp.Packet
}

func listen(c chan message, peer dhcp.Peer) {
	for {
		o, remote, err := peer.Receive(-1)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			continue
		}
		c <- message{remote.IP.String(), o}
	}
}

func snoop(iface string) {

	var mac string
	if iface != "" {
		var err error
		mac, err = MACFromIface(iface)
		checkError(err)
		fmt.Printf("Interface: %s [%s]\n", iface, mac)
	}

	// Set up client
	client, err := dhcp.NewClient()
	checkError(err)
	defer client.Close()

	// Set up server
	server, err := dhcp.NewServer()
	checkError(err)
	defer server.Close()

	c := make(chan message, 1)
	go listen(c, client)
	go listen(c, server)

	for {
		msg := <-c
		stats.pkrec++
		p := msg.packet

		rip := msg.origin
		pmac := p.Chaddr.MACAddress().String()

		var rmac string
		switch rip {
		case "0.0.0.0":
			rmac = pmac
		/*case myip:	// FIXME: check against local ifaces
		rmac = mac*/
		default:
			rmac = MACFromIP(rip)
		}

		if iface != "" && mac != pmac {
			continue
		}

		stats.pkproc++
		stats.count[rmac]++

		if rip == "0.0.0.0" {
			fmt.Printf("\n<<< Broadcast packet\n")
		} else {
			fmt.Printf("\n<<< Packet from %s (%s)\n",
				rip, NameFromIP(rip))
			fmt.Printf("    MAC address: %s (%s)\n",
				rmac, VendorFromMAC(rmac))
		}

		showPacket(&p)
	}
}
