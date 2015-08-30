package main

import (
	"bytes"
	"encoding/binary"
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
	MSG_TYPE_DISCOVER = 1
	MSG_TYPE_OFFER    = 2
	MSG_TYPE_REQUEST  = 3
	MSG_TYPE_DECLINE  = 4
	MSG_TYPE_ACK      = 5
	MSG_TYPE_NACK     = 6
	MSG_TYPE_RELEASE  = 7
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
	Bootp   [192]byte
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

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
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

func sendUDPPacket(p *DHCPPacket, a string) {
	addr, err := net.ResolveUDPAddr("udp4", a)
	checkError(err)
	conn, err := net.DialUDP("udp4", nil, addr)
	checkError(err)
	defer conn.Close()
	sendPacket(conn, p)
}

func NewDHCPDiscoverPacket() *DHCPPacket {
	p := &DHCPPacket{
		Op:    BOOTREQUEST,
		Htype: HTYPE_ETHERNET,
		Hlen:  6,
		Hops:  0,
		Xid:   rand.Uint32(),
		Secs:  0xffff,
		Flags: FLAG_BROADCAST,
		Magic: MAGIC,
		Options: DHCPOptions{OPTION_MSG_TYPE, 1, MSG_TYPE_DISCOVER,
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

func main() {
	iface := "wlp3s0"
	mac := ""
	timeout := 5 * time.Second

	// Get interface MAC address
	ifaces, err := net.Interfaces()
	checkError(err)
	for _, i := range ifaces {
		if i.Name == iface {
			mac = i.HardwareAddr.String()
		}
	}
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
	sendUDPPacket(p, net.IPv4bcast.String() + ":67")

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
