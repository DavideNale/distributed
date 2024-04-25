package main

import (
	"fmt"
	"log"

	"github.com/DavideNale/distributed/p2p"
	"github.com/DavideNale/distributed/storage"
)

type FSOpts struct {
	Root      string
	Transport p2p.Transport
}

type FileServer struct {
	FSOpts
	store  *storage.Store
	quitch chan struct{}
}

func NewFileServer(opts FSOpts) *FileServer {
	return &FileServer{
		FSOpts: opts,
		store:  storage.NewStore(),
		quitch: make(chan struct{}),
	}
}

func (fs *FileServer) Stop() {
	close(fs.quitch)
}

func (fs *FileServer) Start() error {
	if err := fs.Transport.Listen(); err != nil {
		return err
	}

	defer func() {
		log.Println("File server stopped")
		fs.Transport.Close()
	}()

	for {
		select {
		case msg := <-fs.Transport.Consume():
			fmt.Println(msg)
		case <-fs.quitch:
			break
		}
	}
}
