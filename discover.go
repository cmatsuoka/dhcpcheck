package main

import (
	"./dhcp"
	"fmt"
	"net"
	"os"
	"time"
)


func getMAC(s string) (string, error) {
	ifaces, err := net.Interfaces()
	checkError(err)
	for _, i := range ifaces {
		if i.Name == s {
			return i.HardwareAddr.String(), nil
		}
	}
	return "", fmt.Errorf("%s: no such interface", s)
}

func discover(iface string, timeout time.Duration) {

	mac, err := getMAC(iface)
	checkError(err)

	fmt.Printf("Interface: %s [%s]\n", iface, mac)

	// Set up server
	addr, err := net.ResolveUDPAddr("udp4", ":68")
	checkError(err)
	conn, err := net.ListenUDP("udp4", addr)
	checkError(err)
	defer conn.Close()

	// Send discover packet
	p := dhcp.NewDiscoverPacket()
	p.ParseMAC(mac)

	fmt.Println("\n>>> Send DHCP discover")
	showPacket(p)
	err = dhcp.Broadcast(p)
	checkError(err)

	t := time.Now()
	for time.Since(t) < timeout {
		o, remote, err := dhcp.Receive(conn, timeout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			break
		}
		fmt.Println("\n<<< Receive DHCP offer from", remote.IP.String())
		showPacket(&o)
	}
	fmt.Println("No more offers.")
}
