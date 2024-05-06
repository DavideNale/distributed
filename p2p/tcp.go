package p2p

import (
	"errors"
	"fmt"
	"net"

	"github.com/charmbracelet/log"
)

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
	logger     *log.Logger
}

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
		fmt.Println("TCP handshake error")
		return
	}

	msg := &Message{}
	for {
		if err := t.decoder.Decode(conn, msg); err != nil {
			fmt.Println("TCP decode error")
			continue
		}
		fmt.Println(msg)
	}
}
