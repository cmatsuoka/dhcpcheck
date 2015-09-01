package dhcp

import (
	"net"
	"time"
)

func Send(conn net.Conn, p *Packet) error {
	data, err := p.Serialize()
	if err != nil {
		return err
	}

	_, err = conn.Write(data)

	return err
}

func Receive(conn *net.UDPConn, timeout time.Duration) (Packet, *net.UDPAddr, error) {
	var p Packet
	b := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(timeout))
	_, remote, err := conn.ReadFromUDP(b)
	if err != nil {
		return p, remote, err
	}

	err = p.Deserialize(b)

	return p, remote, err
}

func Broadcast(p *Packet) error {
	addr, err := net.ResolveUDPAddr("udp4", net.IPv4bcast.String()+":67")
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	return Send(conn, p)
}
