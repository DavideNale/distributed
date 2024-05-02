package p2p

type HandshakeFunc func(Peer) error

func NoHandshakeFunc(Peer) error { return nil }
