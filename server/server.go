package server

import (
	"github.com/DavideNale/distributed/p2p"
	"github.com/DavideNale/distributed/storage"
	"github.com/charmbracelet/log"
)

type FileServerOpts struct {
	Root           string
	Transport      p2p.Transport
	BootstrapNodes []string
}

type FileServer struct {
	FileServerOpts
	store  *storage.Store
	quitch chan struct{}
	logger *log.Logger
}

func NewFileServer(opts FileServerOpts, logger *log.Logger) *FileServer {
	return &FileServer{
		FileServerOpts: opts,
		store:          storage.NewStore(),
		quitch:         make(chan struct{}),
		logger:         logger,
	}
}

func (s *FileServer) Stop() {
	close(s.quitch)
}

func (s *FileServer) Start() error {
	s.logger.Info("starting fileserver", "port", s.Transport.Addr())
	if err := s.Transport.Listen(); err != nil {
		return err
	}

	defer func() {
		s.logger.Warn("fileserver stopped")
		s.Transport.Close()
	}()

	// Network bootstrapping
	for _, addr := range s.BootstrapNodes {
		if len(addr) == 0 {
			continue
		}
		go func(addr string) {
			s.logger.Debug("attempting bootstrap connection with node", "port", addr)
			if err := s.Transport.Dial(addr); err != nil {
				s.logger.Error("failed to dial node", "port", addr)
			}
		}(addr)
	}

	// for {
	// 	select {
	// 	case msg := <-s.Transport.Consume():
	// 		fmt.Println(msg)
	// 	case <-s.quitch:
	// 		break
	// 	}
	// }
	select {}
	return nil
}
