package main

import (
	"fmt"
	"log"

	"github.com/DavideNale/distributed/p2p"
)

func main() {
	fmt.Println("works")
	tr := p2p.NewTCP(":3000")
	if err := tr.Listen(); err != nil {
		log.Fatal(err)
	}

	select {}
}
