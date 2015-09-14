package main

import (
	"fmt"
	"encoding/json"

	"./dhcp"
	"./format"
	"github.com/cmatsuoka/dncomp"
)

type option struct {
	Len  int
	Name string
}

var (
	options     map[byte]option
	messageType map[byte]string
	op          map[byte]string
)

func init() {
	options = map[byte]option{
		dhcp.PadOption:              {0, "Pad Option"},
		dhcp.Router:                 {-1, "Router"},
		dhcp.SubnetMask:             {4, "Subnet Mask"},
		dhcp.DomainNameServer:       {-1, "Domain Name Server"},
		dhcp.HostName:               {-1, "Host Name"},
		dhcp.DomainName:             {-1, "Domain Name"},
		dhcp.BroadcastAddress:       {4, "Broadcast Address"},
		dhcp.StaticRoute:            {-1, "Static Route"},
		dhcp.IPAddressLeaseTime:     {4, "IP Address Lease Time"},
		dhcp.DHCPMessageType:        {1, "DHCP Message Type"},
		dhcp.ServerIdentifier:       {4, "Server Identifier"},
		dhcp.InterfaceMTU:           {2, "Server MTU"},
		dhcp.RenewalTimeValue:       {4, "Renewal Time Value"},
		dhcp.RebindingTimeValue:     {4, "Rebinding Time Value"},
		dhcp.VendorSpecific:         {-1, "Vendor Specific"},
		dhcp.PerformRouterDiscovery: {1, "Perform Router Discovery"},
		dhcp.NetBIOSNameServer:      {-1, "NetBIOS Name Server"},
		dhcp.NetBIOSNodeType:        {1, "NetBIOS Node Type"},
		dhcp.NetBIOSScope:           {-1, "NetBIOS Scope"},
		dhcp.RequestedIPAddress:     {-1, "Requested IP Address"},
		dhcp.VendorClassIdentifier:  {-1, "Vendor Class Identifier"},
		dhcp.MaxDHCPMessageSize:     {2, "Max DHCP Message Size"},
		dhcp.ParameterRequestList:   {-1, "Parameter Request List"},
		dhcp.ClientIdentifier:       {-1, "Client Identifier"},
		dhcp.DomainSearch:           {-1, "Domain Search"},
		dhcp.UserClass:              {-1, "User Class"},
		dhcp.ClientFQDN:             {-1, "Client FQDN"},
		dhcp.WebProxyServer:         {-1, "Web Proxy Server"},
	}

	messageType = map[byte]string{
		dhcp.DHCPDiscover: "DHCPDISCOVER",
		dhcp.DHCPOffer:    "DHCPOFFER",
		dhcp.DHCPRequest:  "DHCPREQUEST",
		dhcp.DHCPDecline:  "DHCPDECLINE",
		dhcp.DHCPAck:      "DHCPACK",
		dhcp.DHCPNack:     "DHCPNACK",
		dhcp.DHCPRelease:  "DHCPRELEASE",
		dhcp.DHCPInform:   "DHCPINFORM",
	}

	op = map[byte]string{
		dhcp.BootRequest: "BOOTREQUEST",
		dhcp.BootReply:   "BOOTREPLY",
	}
}

func showOptions(p *dhcp.Packet) {

	opts, err := p.DecodeOptions()
	if err != nil {
		fmt.Println("Warning: corrupt option data")
	}

	fmt.Println("Options:")
loop:
	for _, o := range opts {

		switch o.Type {
		case dhcp.EndOption:
			fmt.Print("End Option")
			break loop
		case dhcp.PadOption:
			continue
		}

		if name := options[o.Type].Name; name != "" {
			fmt.Printf("%24s : ", name)
		} else {
			fmt.Printf("%24d : ", o)
		}

		if o.Type == dhcp.VendorClassIdentifier {
			stats.vdc[format.String(o.Data)]++
		}

		switch o.Type {
		case dhcp.DHCPMessageType:
			if m, ok := messageType[o.Data[0]]; ok {
				fmt.Printf(m)
				stats.msg[m]++
			} else {
				s := fmt.Sprintf("<unknown: %d>", o.Data[0])
				fmt.Printf(s)
				stats.msg[s]++
			}

		case dhcp.Router, dhcp.DomainNameServer, dhcp.NetBIOSNameServer:
			// Multiple IP addresses
			for n := 0; n < len(o.Data); n += 4 {
				if n > 0 {
					fmt.Print(", ")
				}
				fmt.Print(format.IPv4String(o.Data[n : n+4]))
			}

		case dhcp.ServerIdentifier, dhcp.SubnetMask,
			dhcp.BroadcastAddress, dhcp.RequestedIPAddress:
			// Single IP address
			fmt.Print(format.IPv4String(o.Data))

		case dhcp.PerformRouterDiscovery:
			// yes or no
			fmt.Print(format.YesNo(o.Data))

		case dhcp.NetBIOSNodeType:
			// hex byte
			fmt.Printf("%#02x", o.Data[0])

		case dhcp.MaxDHCPMessageSize, dhcp.InterfaceMTU:
			// 16-bit integer
			fmt.Print(format.Uint16B(o.Data))

		case dhcp.IPAddressLeaseTime, dhcp.RenewalTimeValue,
			dhcp.RebindingTimeValue:
			// Duration
			fmt.Printf("%d (%s)", format.Uint32B(o.Data),
				format.DurationString(o.Data))

		case dhcp.HostName, dhcp.DomainName, dhcp.WebProxyServer,
			dhcp.NetBIOSScope:
			// String
			fmt.Printf(format.String(o.Data))

		case dhcp.DomainSearch:
			// Compressed domain names (RFC 1035)
			if s, err := dncomp.Decode(o.Data); err != nil {
				fmt.Print(s)
			}

		case dhcp.ClientIdentifier:
			// Types according to RFC 1700
			switch o.Data[0] {
			case 1:
				fmt.Print(format.MACAddrString(o.Data[1:7]))
			default:
				fmt.Printf("type %d (len %d)", o.Data[0], len(o.Data)-1)
			}

		case dhcp.VendorSpecific, dhcp.VendorClassIdentifier, dhcp.UserClass:
			// Dump data
			fmt.Printf(format.String(o.Data))

		case dhcp.ParameterRequestList:
			// Parameter list
			for i, p := range o.Data {
				if i > 0 {
					fmt.Printf("\n%24s   ", "")
				}
				fmt.Printf("%3d %s", p, options[p].Name)
			}

		case dhcp.ClientFQDN:
			// Client FQDN format
			c := []byte{'-', '-', '-', '-'}
			d := []byte{'N', 'E', 'O', 'S'}
			for j := range c {
				if o.Data[0]&(1<<(3-uint(j))) != 0 {
					c[j] = d[j]
				}
			}
			fmt.Printf("%s %02x %02x ", string(c), o.Data[1],
				o.Data[2])
			if o.Data[0]&0x04 == 0 {
				fmt.Printf("%q", string(o.Data[3:]))
			} else {
				fmt.Printf("%q", format.CanonicalWireFormat(o.Data[3:]))
			}
		}
		fmt.Println()
	}
}

func opcode(o byte) string {
	if s := op[o]; s != "" {
		return s
	}
	return fmt.Sprintf("<unknown:%d>", o)
}

func showPacket(p *dhcp.Packet) {
	fmt.Printf("Message opcode    : %s\n", opcode(p.Op))
	//fmt.Printf("HW address type   : %d\n", p.Htype)
	//fmt.Printf("HW address length : %d\n", p.Hlen)
	//fmt.Printf("Hops              : %d\n", p.Hops)
	fmt.Printf("Transaction ID    : %#08x\n", p.Xid)
	//fmt.Printf("Seconds elapsed   : %d\n", p.Secs)
	fmt.Printf("Flags             : %#04x\n", p.Flags)
	fmt.Printf("Client IP address : %s\n", p.Ciaddr.String())
	fmt.Printf("Your IP address   : %s\n", p.Yiaddr.String())
	fmt.Printf("Server IP address : %s\n", p.Siaddr.String())
	fmt.Printf("Relay IP address  : %s\n", p.Giaddr.String())

	mac := p.Chaddr.MACAddress().String()
	fmt.Printf("Client MAC address: %s (%s)\n", mac, VendorFromMAC(mac))

	showOptions(p)

	fmt.Println()

	// Update report

	report.Packets++
	report.MsgType = stats.msg

        vcount := map[string]uint{}
        for key, val := range stats.count {
                v := VendorFromMAC(key)
                vcount[v] += val
        }
	report.Vendors = vcount
	report.VdClass = stats.vdc

	j,err := json.Marshal(report)
	if err != nil {
		fmt.Errorf("Error: %s\n", err.Error())
		return
	}
	repch <- string(j)
}
