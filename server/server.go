package server

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"sync"

	"github.com/DavideNale/distributed/p2p"
	"github.com/DavideNale/distributed/storage"
	"github.com/charmbracelet/log"
)

// FileServerOpts is the options for the FileServer.
type FileServerOpts struct {
	Root           string
	Transport      p2p.Transport
	BootstrapNodes []string
	Logger         *log.Logger
}

// FileServer is a FileServer that contains a Store.
type FileServer struct {
	FileServerOpts

	store  *storage.Store
	quitch chan struct{}

	peerLock sync.Mutex
	peers    map[string]p2p.Peer
}

// NewFileServer returns a pointer to a FileServer with the specified options.
func NewFileServer(opts FileServerOpts, logger *log.Logger) *FileServer {
	storeOpts := storage.StoreOpts{
		Root:            opts.Root,
		PathTransformer: storage.HashTransformer,
	}
	server := &FileServer{
		FileServerOpts: opts,
		store:          storage.NewStore(storeOpts),
		quitch:         make(chan struct{}),
		peerLock:       sync.Mutex{},
		peers:          make(map[string]p2p.Peer),
	}

	server.Transport.(*p2p.TCP).OnPeer = server.OnPeer
	return server
}

// Stop sends a stopping signal to the server via the quitch channel.
func (s *FileServer) Stop() {
	close(s.quitch)
}

// Start bootstraps the server and starts the listen loop.
func (s *FileServer) Start() error {
	s.Transport.(*p2p.TCP).OnPeer = s.OnPeer // TODO: move

	s.Logger.Info("starting fileserver", "port", s.Transport.Addr())
	if err := s.Transport.Listen(); err != nil {
		return err
	}

	defer func() {
		s.Logger.Warn("fileserver stopped")
		s.Transport.Close()
	}()

	// Network bootstrapping
	for _, addr := range s.BootstrapNodes {
		if len(addr) == 0 {
			continue
		}
		go func(addr string) {
			s.Logger.Debug("attempting bootstrap connection with node", "port", addr)
			if err := s.Transport.Dial(addr); err != nil {
				s.Logger.Error("failed to dial node", "port", addr)
			}
		}(addr)
	}

	// Loop
	for {
		select {
		case rpc := <-s.Transport.Consume():
			var msg Message
			if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg); err != nil {
				fmt.Printf("%+v\n", err)
				s.Logger.Error("decoding error")
			}
			if err := s.handleMessage(rpc.From, &msg); err != nil {
				s.Logger.Error("error handling message", "from", rpc.From)
			}
		// 	fmt.Println(msg)
		case <-s.quitch:
			return nil
		}
	}
}

// OnPeer adds a peer to the peers list of the FilesServer
func (s *FileServer) OnPeer(p p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	addr := p.RemoteAddr().String()
	s.peers[addr] = p
	s.Logger.Info("added a new peer", "address", addr)
	return nil
}
