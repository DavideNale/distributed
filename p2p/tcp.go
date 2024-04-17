package p2p

import (
	"fmt"
	"net"
)

// Peer represents a remote node in the network.
type Peer interface{}

// Transport represents anything that handles the
// communication between nodes.
type Transport interface {
	Listen() error
}

type TCPPeer struct {
	conn     net.Conn
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCP struct {
	listenAddr string
	listener   net.Listener
	handshake  HandshakeFunc
	decoder    Decoder
}

func NewTCP(listenAddr string) *TCP {
	return &TCP{
		listenAddr: listenAddr,
		handshake:  NOPHandshakeFunc,
	}
}

func (t *TCP) Listen() error {
	var err error
	t.listener, err = net.Listen("tcp", t.listenAddr)
	if err != nil {
		return err
	}

	go func() {
		for {
			conn, err := t.listener.Accept()
			if err != nil {
				fmt.Println("TCP error")
			}
			fmt.Println("TCP connection accepted")
			go t.handleConn(conn)
		}
	}()

	return nil
}

func (t *TCP) handleConn(conn net.Conn) {
	// peer := NewTCPPeer(conn, true)

	if err := t.handshake(conn); err != nil {
		conn.Close()
		fmt.Println("TCP handshake error")
		return
	}

	buf := make([]byte, 1000)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("TCP decode error")
		}
		fmt.Println(buf[:n])
	}
}
