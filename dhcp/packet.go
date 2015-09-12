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
	DHCPInform   = 8
)

const (
	FlagBroadcast = 1 << 15
)

const (
	packetSize    = 548
	magic         = 0x63825363
	HtypeEthernet = 1
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type IPv4Address [4]byte

func (a *IPv4Address) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", a[0], a[1], a[2], a[3])
}

type MACAddress [6]byte

func (a *MACAddress) String() string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
		a[0], a[1], a[2], a[3], a[4], a[5])
}

type HWAddress [16]byte

func (a *HWAddress) MACAddress() *MACAddress {
	var mac MACAddress
	copy(mac[:], a[:6])
	return &mac
}

type OptionsArea [1200]byte

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
	Chaddr  HWAddress   // client hardware address
	Sname   [64]byte
	File    [128]byte
	Magic   uint32
	Options OptionsArea
}

func (p Packet) serialize() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, p)
	return buf.Bytes(), err
}

func (p *Packet) deserialize(data []byte) error {
	return binary.Read(bytes.NewReader(data), binary.BigEndian, p)
}

// SetClientMAC takes a MAC address and sets the client hardware address
// field of the DHCP packet.
func (p *Packet) SetClientMAC(mac string) error {
	hw, err := net.ParseMAC(mac)
	if err == nil {
		copy(p.Chaddr[0:6], hw)
	}
	return err
}

// NewDiscoverPacket builds a new DHCPDISCOVER packet.
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

func (p *Packet) DecodeOptions() ([]Option, error) {

	var option []Option

	for i := 0; i < len(p.Options); {

		if i >= len(p.Options) {
			return option, ErrCorruptedOptions
		}

		o := p.Options[i]
		i++

		if o == EndOption || o == PadOption {
			option = append(option, Option{o, nil})
			continue
		}

		if i >= len(p.Options) {
			return option, ErrCorruptedOptions
		}

		l := int(p.Options[i])
		i++

		if i+l > len(p.Options) {
			return option, ErrCorruptedOptions
		}

		option = append(option, Option{o, p.Options[i : i+l]})

		i += l
	}

	return option, nil
}
