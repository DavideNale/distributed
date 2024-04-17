package main

import (
	"github.com/DavideNale/distributed/p2p"
)

func main() {
	tr := p2p.NewTCP(":3000")
	tr.Listen()

	select {}
}
