package main

import (
	"./dhcp"
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/cmatsuoka/dncomp"
)

type option struct {
	Len  int
	Name string
}

var options map[byte]option
var messageType map[byte]string

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
		dhcp.VendorSpecific:     {-1, "Vendor Specific"},
		dhcp.NetBIOSNameServer:  {-1, "NetBIOS Name Server"},
		dhcp.DomainSearch:       {-1, "Domain Search"},
		dhcp.WebProxyServer:     {-1, "Web Proxy Server"},
	}

	messageType = map[byte]string{
		dhcp.DHCPDiscover: "DHCPDISCOVER",
		dhcp.DHCPOffer:    "DHCPOFFER",
		dhcp.DHCPRequest:  "DHCPREQUEST",
		dhcp.DHCPDecline:  "DHCPDECLINE",
		dhcp.DHCPAck:      "DHCPACK",
		dhcp.DHCPNack:     "DHCPNACK",
		dhcp.DHCPRelease:  "DHCPRELEASE",
	}
}

func showOptions(p *dhcp.Packet) {
	b32 := func(data []byte) uint32 {
		buf := bytes.NewBuffer(data)
		var x uint32
		binary.Read(buf, binary.BigEndian, &x)
		return x
	}

	ip4 := func(data []byte) string {
		var ip dhcp.IPv4Address
		copy(ip[:], data[0:4])
		return ip.String()
	}

	opts := p.Options
	fmt.Println("Options:")
loop:
	for i := 0; i < len(opts); {
		o := opts[i]
		i++

		switch o {
		case dhcp.EndOption:
			fmt.Print("End Option")
			break loop
		case dhcp.PadOption:
			continue
		}

		length := int(opts[i])
		i++
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
			fmt.Print(messageType[opts[i]])

		case dhcp.Router, dhcp.DomainNameServer, dhcp.NetBIOSNameServer:
			// Multiple IP addresses
			for n := 0; n < length; n += 4 {
				fmt.Print(ip4(opts[i+n:i+4+n]), " ")
			}

		case dhcp.ServerIdentifier, dhcp.SubnetMask, dhcp.BroadcastAddress:
			// Single IP address
			fmt.Print(ip4(opts[i:]))

		case dhcp.IPAddressLeaseTime, dhcp.RenewalTimeValue, dhcp.RebindingTimeValue:
			// 32-bit integer
			fmt.Print(b32(opts[i:]))

		case dhcp.HostName, dhcp.DomainName, dhcp.WebProxyServer:
			// String
			fmt.Print(string(opts[i : i+length]))

		case dhcp.DomainSearch:
			// Compressed domain names (RFC 1035)
			if s, err := dncomp.Decode(opts[i : i+length]); err != nil {
				fmt.Print(s)
			}

		case dhcp.VendorSpecific:
			// Size only
			fmt.Printf("(%d bytes)", length)
		}
		fmt.Println()

		i += length
	}
}

func showPacket(p *dhcp.Packet) {
	fmt.Println("Client IP address :", p.Ciaddr.String())
	fmt.Println("Your IP address   :", p.Yiaddr.String())
	fmt.Println("Server IP address :", p.Siaddr.String())
	fmt.Println("Relay IP address  :", p.Giaddr.String())
	showOptions(p)
	fmt.Println()
}
