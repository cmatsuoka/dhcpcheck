package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
)

const (
	BOOTREQUEST = 1
	BOOTREPLY   = 2
)

const (
	OPTION_MSG_TYPE = 0x35
	OPTION_END      = 0xff
)

const (
	DHCPDISCOVER = 1
	DHCPOFFER    = 2
	DHCPREQUEST  = 3
	DHCPDECLINE  = 4
	DHCPACK      = 5
	DHCPNACK     = 6
	DHCPRELEASE  = 7
)

const (
	FLAG_BROADCAST = 1 << 15
)

const (
	PACKET_SIZE    = 548
	HTYPE_ETHERNET = 1
	MAGIC          = 0x63825363
)

type IPv4Address [4]byte

func (a *IPv4Address) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", a[0], a[1], a[2], a[3])
}

type DHCPOptions [308]byte

type DHCPPacket struct {
	Op      byte
	Htype   byte
	Hlen    byte
	Hops    byte
	Xid     uint32
	Secs    uint16
	Flags   uint16
	Ciaddr  IPv4Address // client IP address
	Yiaddr  IPv4Address // your IP address
	Siaddr  IPv4Address // server IP address
	Giaddr  IPv4Address // gateway IP address
	Chaddr  [16]byte    // client hardware address
	Sname   [64]byte
	File    [128]byte
	Magic   uint32
	Options DHCPOptions
}

func (p *DHCPPacket) parseMAC(s string) error {
	hw, err := net.ParseMAC(s)
	if err == nil {
		copy(p.Chaddr[0:6], hw)
	}
	return err
}

func sendPacket(conn net.Conn, p *DHCPPacket) error {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, *p)
	if err != nil {
		return err
	}
	_, err = conn.Write(buffer.Bytes())
	return err
}

func receivePacket(conn *net.UDPConn, timeout time.Duration) (*DHCPPacket, *net.UDPAddr, error) {
	conn.SetReadDeadline(time.Now().Add(timeout))
	var b [1024]byte
	_, remote, err := conn.ReadFromUDP(b[:])
	if err != nil {
		return nil, remote, err
	}
	var p DHCPPacket
	err = binary.Read(bytes.NewReader(b[:]), binary.BigEndian, &p)
	return &p, remote, err
}

func sendUDPPacket(p *DHCPPacket, a string) error {
	addr, err := net.ResolveUDPAddr("udp4", a)
	if err != nil {
		return err
	}
	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	return sendPacket(conn, p)
}

func NewDHCPDiscoverPacket() *DHCPPacket {
	p := &DHCPPacket{
		Op:    BOOTREQUEST,
		Htype: HTYPE_ETHERNET,
		Hlen:  6,
		Hops:  0,
		Xid:   rand.Uint32(),
		Secs:  0,
		Flags: FLAG_BROADCAST,
		Magic: MAGIC,
		Options: DHCPOptions{OPTION_MSG_TYPE, 1, DHCPDISCOVER,
			OPTION_END},
	}

	return p
}

func showOffer(p *DHCPPacket) {
	fmt.Println("Client IP address :", p.Ciaddr.String())
	fmt.Println("Your IP address   :", p.Yiaddr.String())
	fmt.Println("Server IP address :", p.Siaddr.String())
	fmt.Println("Relay IP address  :", p.Giaddr.String())
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
	p := NewDHCPDiscoverPacket()
	p.parseMAC(mac)
	err = sendUDPPacket(p, net.IPv4bcast.String()+":67")
	checkError(err)

	t := time.Now()
	for time.Since(t) < timeout {
		o, remote, err := receivePacket(conn, timeout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			break
		}
		fmt.Println("Receive DHCP offer from", remote.IP.String())
		showOffer(o)
	}
	fmt.Println("No more offers.")
}
