package dhcp

import (
	"fmt"
	"net"
	"time"
)

// Peer definitions

type peer struct {
	local      *net.UDPConn
	remote     net.Conn
	localPort  int
	remotePort int
}

func newPeer(localPort, remotePort int, listen bool) (*peer, error) {

	pr := &peer{}
	pr.localPort = localPort
	pr.remotePort = remotePort

	if !listen {
		return pr, nil
	}

	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf(":%d", localPort))
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		return nil, err
	}

	pr.local = conn

	return pr, nil
}

func (pr *peer) setRemote(ip net.IP) error {
	addr, err := net.ResolveUDPAddr("udp4",
		fmt.Sprintf("%s:%d", ip.String(), pr.remotePort))
	if err != nil {
		return err
	}
	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		return err
	}

	pr.remote = conn

	return nil
}

func (pr *peer) Close() {
	pr.local.Close()
}

func (pr *peer) closeRemote() {
	pr.remote.Close()
}

func (pr *peer) Send(p *Packet) error {
	return send(pr.remote, p)
}

func (pr *peer) Receive(timeout time.Duration) (Packet, *net.UDPAddr, error) {
	return receive(pr.local, timeout)
}

func (pr *peer) Broadcast(p *Packet) error {
	addr, err := net.ResolveUDPAddr("udp4",
		fmt.Sprintf("%s:%d", net.IPv4bcast.String(), pr.remotePort))
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	return send(conn, p)
}

// Client

type Client struct {
	peer
}

func NewClient() (*Client, error) {
	pr, err := newPeer(68, 67, true)
	if pr == nil {
		return nil, err
	}
	return &Client{*pr}, err
}

func NewClientNotListening() (*Client, error) {
	pr, err := newPeer(68, 67, false)
	if pr == nil {
		return nil, err
	}
	return &Client{*pr}, err
}

func (cl *Client) SetServer(svIP net.IP) error {
	return cl.setRemote(svIP)
}

func (cl *Client) CloseServer() {
	cl.closeRemote()
}

// Server

type Server struct {
	peer
}

func NewServer() (*Server, error) {
	pr, err := newPeer(67, 68, true)
	if pr == nil {
		return nil, err
	}
	return &Server{*pr}, err
}

func (sv *Server) SetClient(clIP net.IP) error {
	return sv.setRemote(clIP)
}

func (sv *Server) CloseClient() {
	sv.closeRemote()
}

// Helpers

func send(conn net.Conn, p *Packet) error {
	data, err := p.serialize()
	if err != nil {
		return err
	}

	_, err = conn.Write(data)

	return err
}

func receive(conn *net.UDPConn, timeout time.Duration) (Packet, *net.UDPAddr, error) {
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
