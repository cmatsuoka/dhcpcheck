package dhcp

import (
	"math/rand"
	"net"
)

type DiscoverPacket struct {
	Packet
}

func (p *DiscoverPacket) Send() error {
	addr, err := net.ResolveUDPAddr("udp4", net.IPv4bcast.String()+":67")
	if err != nil {
		return err
	}
	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	return p.Packet.Send(conn)
}

func NewDiscoverPacket() *DiscoverPacket {
	p := &DiscoverPacket{
		Packet{
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
		},
	}

	return p
}
