package main

import (
	"./dhcp"
	"flag"
	"fmt"
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

	mac, err := MACFromIface(iface)
	checkError(err)

	fmt.Printf("Interface: %s [%s]\n", iface, mac)

	var client *dhcp.Client

	if timeout <= 0 {
		client, err = dhcp.NewClientNotListening()
		checkError(err)
	} else {
		client, err = dhcp.NewClient()
		checkError(err)
		defer client.Close()
	}

	// Send discover packet
	p := dhcp.NewDiscoverPacket()
	p.SetClientMAC(mac)

	fmt.Println("\n>>> Send DHCP discover")
	showPacket(p)
	err = client.Broadcast(p)
	checkError(err)

	if timeout <= 0 {
		return
	}

	t := time.Now()
	for time.Since(t) < timeout {
		o, remote, err := client.Receive(timeout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			break
		}

		if o.Xid == p.Xid {
			ip := remote.IP.String()
			mac := MACFromIP(ip)

			fmt.Printf("\n<<< Receive DHCP offer from %s (%s)\n",
				ip, NameFromIP(ip))
			fmt.Printf("    MAC address: %s (%s)\n",
				mac, VendorFromMAC(mac))

			showPacket(&o)
		}
	}
	fmt.Println("No more offers.")
}
