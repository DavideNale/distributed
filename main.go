package main

import (
	"time"

	"github.com/DavideNale/distributed/p2p"
	"github.com/DavideNale/distributed/server"
	"github.com/charmbracelet/log"
)

func main() {
	s1 := makeServer(":5000", "s1/")
	s2 := makeServer(":5001", "s2/", ":5000")

	go func() { log.Fatal(s1.Start()) }()
	time.Sleep(2 * time.Millisecond)
	go func() { log.Fatal(s2.Start()) }()

	select {}
}

// makeServer initializes a FileServer and returns a pointer to it
func makeServer(port string, root string, nodes ...string) *server.FileServer {
	tcpTransport := p2p.NewTCP(port)
	fileServerOpts := server.FileServerOpts{
		Root:           root,
		Transport:      tcpTransport,
		BootstrapNodes: nodes,
	}
	return server.NewFileServer(fileServerOpts)
}
