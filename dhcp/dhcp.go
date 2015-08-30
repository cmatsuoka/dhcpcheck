package dhcp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"time"
)

const (
	BOOTREQUEST = 1
	BOOTREPLY   = 2
)

const (
	DHCPDiscover = 1
	DHCPOffer    = 2
	DHCPRequest  = 3
	DHCPDecline  = 4
	DHCPAck      = 5
	DHCPNack     = 6
	DHCPRelease  = 7
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

type OptionsArea [308]byte

type Packet struct {
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
	Options OptionsArea
}

func (p *Packet) ParseMAC(s string) error {
	hw, err := net.ParseMAC(s)
	if err == nil {
		copy(p.Chaddr[0:6], hw)
	}
	return err
}

func SendPacket(conn net.Conn, p *Packet) error {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, *p)
	if err != nil {
		return err
	}
	_, err = conn.Write(buffer.Bytes())
	return err
}

func ReceivePacket(conn *net.UDPConn, timeout time.Duration) (*Packet, *net.UDPAddr, error) {
	conn.SetReadDeadline(time.Now().Add(timeout))
	var b [1024]byte
	_, remote, err := conn.ReadFromUDP(b[:])
	if err != nil {
		return nil, remote, err
	}
	var p Packet
	err = binary.Read(bytes.NewReader(b[:]), binary.BigEndian, &p)
	return &p, remote, err
}

func SendUDPPacket(p *Packet, a string) error {
	addr, err := net.ResolveUDPAddr("udp4", a)
	if err != nil {
		return err
	}
	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	return SendPacket(conn, p)
}

func DiscoverPacket() *Packet {
	p := &Packet{
		Op:    BOOTREQUEST,
		Htype: HTYPE_ETHERNET,
		Hlen:  6,
		Hops:  0,
		Xid:   rand.Uint32(),
		Secs:  0,
		Flags: FLAG_BROADCAST,
		Magic: MAGIC,
		Options: OptionsArea{DHCPMessageType, 1, DHCPDiscover,
			EndOption},
	}

	return p
}
