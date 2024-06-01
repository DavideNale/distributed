package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"time"

	"github.com/DavideNale/distributed/p2p"
	"github.com/DavideNale/distributed/server"
	"github.com/charmbracelet/log"
)

func main() {
	s1 := makeServer(":5000", "s1/")
	s2 := makeServer(":5001", "s2/", ":5000")
	s3 := makeServer(":5002", "s3/", ":5001")

	go func() { log.Fatal(s1.Start()) }()
	time.Sleep(2 * time.Millisecond)
	go func() { log.Fatal(s2.Start()) }()
	time.Sleep(2 * time.Millisecond)
	go func() { log.Fatal(s3.Start()) }()

	time.Sleep(1 * time.Second)

	key := "file.png"
	data := bytes.NewReader([]byte("file content"))

	s1.Store(key, data)

	time.Sleep(3 * time.Second)
	reader, _ := s2.Get(key)
	content := readFileContent(reader)
	s2.Logger.Info(key, "content", content)

	time.Sleep(3 * time.Second)
	s3.Delete(key)
	time.Sleep(2 * time.Second)
	// s1.Clear()
	// s2.Clear()
	// s3.Clear()
	select {}
}

// makeServer returns a pointer to a configured server.
func makeServer(port string, root string, nodes ...string) *server.FileServer {
	logger := log.NewWithOptions(os.Stderr, log.Options{
		Level:  log.DebugLevel,
		Prefix: strings.Split(root, "/")[0],
	})
	fileServerOpts := server.FileServerOpts{
		Root:           root,
		Transport:      p2p.NewTCP(port, logger),
		BootstrapNodes: nodes,
		Logger:         logger,
	}
	return server.NewFileServer(fileServerOpts, logger)
}

// readFileContent reads the file content as a string.
func readFileContent(r io.Reader) string {
	file, _ := io.ReadAll(r)
	return string(file)
}
