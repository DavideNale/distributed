package p2p

import (
	"errors"
	"fmt"
	"net"

	"github.com/charmbracelet/log"
)

// TCPPeer representa a node over an established TCP connection.
type TCPPeer struct {
	conn     net.Conn
	outbound bool
}

// NewTCPPeer returns a pointer to a TCPPeer.
func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

// TCP implements the Trasport interface.
type TCP struct {
	listenAddr string
	listener   net.Listener
	handshake  HandshakeFunc
	decoder    Decoder
	logger     *log.Logger
}

// NewTCP returns a pointer to a newly configured TPC.
func NewTCP(listenAddr string, logger *log.Logger) *TCP {
	return &TCP{
		listenAddr: listenAddr,
		handshake:  NoHandshakeFunc,
		decoder:    DefaultDecoder{},
		logger:     logger,
	}
}

// Close implements the Transport interface.
func (t *TCP) Close() error {
	return t.listener.Close()
}

// Addr implements the Transport interface.
func (t *TCP) Addr() string {
	return t.listenAddr
}

// Dial implements the Transport interface.
func (t *TCP) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	go t.handleConn(conn)
	return nil
}

// Listen implements the Transport interface.
func (t *TCP) Listen() error {
	var err error
	t.listener, err = net.Listen("tcp", t.listenAddr)
	if err != nil {
		return err
	}

	go func() {
		for {
			conn, err := t.listener.Accept()
			if errors.Is(err, net.ErrClosed) {
				return
			}
			if err != nil {
				fmt.Println("TCP error")
			}
			t.logger.Debug("TCP connection accepted", "port", t.listenAddr)
			go t.handleConn(conn)
		}
	}()

	return nil
}

func (t *TCP) handleConn(conn net.Conn) {

	if err := t.handshake(conn); err != nil {
		conn.Close()
		t.logger.Error("TCP handshake error")
		return
	}

	msg := &Message{}
	for {
		if err := t.decoder.Decode(conn, msg); err != nil {
			t.logger.Error("TCP error decoding RPC")
			continue
		}
		fmt.Println(msg)
	}
}

func (t *TCP) Consume() <-chan RPC {
	return make(chan RPC, 1014)
}
