package p2p

// Peer represents a remote node in the network.
type Peer interface{}

// Transport represents the communication channel.
type Transport interface {
	Addr() string
	Listen() error
	Dial(string) error
	Close() error
}
