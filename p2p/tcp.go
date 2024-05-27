package p2p

import (
	"errors"
	"net"
	"sync"

	"github.com/charmbracelet/log"
)

// TCPPeer representa a node over an established TCP connection.
type TCPPeer struct {
	net.Conn
	outbound bool
	wg       *sync.WaitGroup
}

// NewTCPPeer returns a pointer to a TCPPeer.
func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn:     conn,
		outbound: outbound,
		wg:       &sync.WaitGroup{},
	}
}

func (p *TCPPeer) CloseStream() {
	p.wg.Done()
}

func (p *TCPPeer) Send(b []byte) error {
	_, err := p.Conn.Write(b)
	return err
}

// TCP implements the Trasport interface.
type TCP struct {
	listenAddr string
	listener   net.Listener
	handshake  HandshakeFunc
	decoder    Decoder
	logger     *log.Logger
	rpcch      chan RPC
	OnPeer     func(Peer) error
}

// NewTCP returns a pointer to a newly configured TPC.
func NewTCP(listenAddr string, logger *log.Logger) *TCP {
	return &TCP{
		listenAddr: listenAddr,
		handshake:  NoHandshakeFunc,
		decoder:    DefaultDecoder{},
		logger:     logger,
		rpcch:      make(chan RPC, 1024),
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
	go t.handleConn(conn, false)
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
				t.logger.Error("TCP error")
			}
			t.logger.Debug("TCP connection accepted", "port", t.listenAddr)
			go t.handleConn(conn, false)
		}
	}()

	return nil
}

func (t *TCP) handleConn(conn net.Conn, outbound bool) {
	defer func() {
		t.logger.Debug("dropping peer connection", "peer", conn.RemoteAddr().String())
		defer conn.Close()
	}()

	peer := NewTCPPeer(conn, outbound)
	if err := t.handshake(peer); err != nil {
		return
	}
	if t.OnPeer != nil {
		if err := t.OnPeer(peer); err != nil {
			return
		}
	}

	rpc := RPC{}
	for {
		if err := t.decoder.Decode(conn, &rpc); err != nil {
			t.logger.Error("TCP error decoding RPC")
			return
		}
		rpc.From = conn.RemoteAddr().String()
		t.rpcch <- rpc
		peer.wg.Add(1)
		peer.wg.Wait()
	}
}

func (t *TCP) Consume() <-chan RPC {
	return t.rpcch
}
