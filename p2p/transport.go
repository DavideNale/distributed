package p2p

// Peer represents a remote node in the network.
type Peer interface{}

// Transport represents the communication channel.
type Transport interface {
	Addr() string
	Listen() error
	Dial(string) error
	Consume() <-chan RPC
	Close() error
}

// RPC represents an arbitrary message sent over a Transport.
type RPC struct {
	From    string
	Payload []byte
	Stream  bool
}
