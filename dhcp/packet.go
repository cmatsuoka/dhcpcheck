package dhcp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
)

const (
	BootRequest = 1
	BootReply   = 2
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
	FlagBroadcast = 1 << 15
)

const (
	packetSize    = 548
	magic         = 0x63825363
	HtypeEthernet = 1
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

func (p Packet) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, p)
	return buf.Bytes(), err
}

func (p *Packet) Deserialize(data []byte) error {
	return binary.Read(bytes.NewReader(data), binary.BigEndian, p)
}

func NewDiscoverPacket() *Packet {
	p := &Packet{
		Op:    BootRequest,
		Htype: HtypeEthernet,
		Hlen:  6,
		Hops:  0,
		Xid:   rand.Uint32(),
		Secs:  0,
		Flags: FlagBroadcast,
		Magic: magic,
		Options: OptionsArea{DHCPMessageType, 1, DHCPDiscover,
			EndOption},
	}

	return p
}
