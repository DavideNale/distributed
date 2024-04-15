package p2p

import (
	"fmt"
	"net"
)

type Peer interface{}

type TCPPeer struct {
	conn     net.Conn
	outbound bool
}

type TCP struct {
	listenAddr string
	listener   net.Listener
}

func NewTCP(listenAddr string) *TCP {
	return &TCP{
		listenAddr: listenAddr,
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
			fmt.Println("TCP connection accpeted")
			_ = conn
		}
	}()

	return nil
}
