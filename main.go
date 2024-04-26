package main

import (
	"log"

	"github.com/DavideNale/distributed/p2p"
)

func main() {

	tr := p2p.NewTCP(":3000")
	tr.Listen()
	ops := FSOpts{
		Root:      "root_3000",
		Transport: tr,
	}
	s := NewFileServer(ops)

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}

	select {}
}
