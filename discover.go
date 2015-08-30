package main

import (
	"./dhcp"
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

type option struct {
	Len  int
	Name string
}

var options map[byte]option

func init() {
	options = map[byte]option{
		dhcp.PadOption:          {0, "Pad Option"},
		dhcp.Router:             {-1, "Router"},
		dhcp.SubnetMask:         {4, "Subnet Mask"},
		dhcp.DomainNameServer:   {-1, "Domain Name Server"},
		dhcp.HostName:           {-1, "Host Name"},
		dhcp.DomainName:         {-1, "Domain Name"},
		dhcp.BroadcastAddress:   {4, "Broadcast Address"},
		dhcp.StaticRoute:        {-1, "Static Route"},
		dhcp.IPAddressLeaseTime: {4, "IP Address Lease Time"},
		dhcp.DHCPMessageType:    {1, "DHCP Message Type"},
		dhcp.ServerIdentifier:   {4, "Server Identifier"},
		dhcp.RenewalTimeValue:   {4, "Renewal Time Value"},
		dhcp.RebindingTimeValue: {4, "Rebinding Time Value"},
		dhcp.DomainSearch:       {-1, "Domain Search"},
		dhcp.WebProxyServer:     {-1, "Web Proxy Server"},
	}
}

func b32(data []byte) uint32 {
	buf := bytes.NewBuffer(data)
	var x uint32
	binary.Read(buf, binary.BigEndian, &x)
	return x
}

func ip4(data []byte) string {
	var ip dhcp.IPv4Address
	copy(ip[:], data[0:4])
	return ip.String()
}

func parseOptions(p *dhcp.Packet) {
	opts := p.Options
	fmt.Println("Options:")
loop:
	for i := 0; i < len(opts); i++ {
		o := opts[i]

		switch o {
		case dhcp.EndOption:
			fmt.Print("End Option")
			break loop
		case dhcp.PadOption:
			continue
		}

		length := int(opts[i+1])
		_, ok := options[o]
		if ok && options[o].Len >= 0 && options[o].Len != length {
			fmt.Printf("corrupted option (%d,%d)\n",
				options[o].Len, length)
			break loop
		}

		if name := options[o].Name; name != "" {
			fmt.Printf("%24s : ", options[o].Name)
		} else {
			fmt.Printf("%24d : ", o)
		}

		switch o {
		case dhcp.DHCPMessageType:
			t := opts[i+2]
			fmt.Print(t)
			break
		case dhcp.Router, dhcp.DomainNameServer:
			for n := 0; n < length; n += 4  {
				fmt.Print(ip4(opts[i+2+n:i+6+n]), " ")
			}
		case dhcp.ServerIdentifier, dhcp.SubnetMask, dhcp.BroadcastAddress:
			fmt.Print(ip4(opts[1+2:i+6]))
			break
		case dhcp.IPAddressLeaseTime, dhcp.RenewalTimeValue, dhcp.RebindingTimeValue:
			fmt.Print(b32(opts[i+2 : i+6]))
			break
		case dhcp.HostName, dhcp.DomainName, dhcp.WebProxyServer:
			fmt.Print(string(opts[i+2:i+2+length]))
			break
		case dhcp.DomainSearch:
			fmt.Print("[TODO RFC 1035 section 4.1.4]")
			break
		}
		fmt.Println()

		i += 1 + length
	}
}

func showOffer(p *dhcp.Packet) {
	fmt.Println("Client IP address :", p.Ciaddr.String())
	fmt.Println("Your IP address   :", p.Yiaddr.String())
	fmt.Println("Server IP address :", p.Siaddr.String())
	fmt.Println("Relay IP address  :", p.Giaddr.String())
	parseOptions(p)
	fmt.Println()
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}

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

	mac := ""
	timeout := time.Duration(secs) * time.Second

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
	fmt.Println("Send DHCP discover\n")
	p := dhcp.DiscoverPacket()
	p.ParseMAC(mac)
	err = dhcp.SendUDPPacket(p, net.IPv4bcast.String()+":67")
	checkError(err)

	t := time.Now()
	for time.Since(t) < timeout {
		o, remote, err := dhcp.ReceivePacket(conn, timeout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			break
		}
		fmt.Println("Receive DHCP offer from", remote.IP.String())
		showOffer(o)
	}
	fmt.Println("No more offers.")
}
