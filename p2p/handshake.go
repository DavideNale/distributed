package p2p

type HandshakeFunc func(any) error

var NOPHandshakeFunc = func(any) error { return nil }
