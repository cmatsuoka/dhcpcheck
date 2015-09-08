package dhcp

import (
	"net"
	"time"
)

// Send sends a packet using the supplied network connection.
func Send(conn net.Conn, p *Packet) error {
	data, err := p.serialize()
	if err != nil {
		return err
	}

	_, err = conn.Write(data)

	return err
}

// Receive reads a new packet from the supplied network connection,
// optionally waiting for the amount of time specified in timeout.
func Receive(conn *net.UDPConn, timeout time.Duration) (Packet, *net.UDPAddr, error) {
	var p Packet
	b := make([]byte, 1024)
	if timeout > 0 {
		conn.SetReadDeadline(time.Now().Add(timeout))
	}
	_, remote, err := conn.ReadFromUDP(b)
	if err != nil {
		return p, remote, err
	}

	err = p.deserialize(b)

	return p, remote, err
}

// Broadcast broadcasts the DHCP packet to the server.
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
