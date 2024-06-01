# Distributed File System in Go

## Overview

This project implements a distributed file system in Go, utilizing a gossip protocol over TCP streaming to ensure file consistency and availability across multiple peers. The system supports basic operations such as storing, retrieving, deleting, and clearing files.

![](demo.gif)
## Features

- **Distributed Storage**: Files are stored and retrieved across a distributed network of peers.
- **Gossip Protocol**: Changes are propagated using a gossip protocol to ensure eventual consistency.
- **TCP Streaming**: Efficient data transfer using TCP streaming.
- **Logging**: Detailed logging of operations for easier debugging and monitoring.

## Example

To run this project:
1. Clone the repo:
``` bash
git clone https://github.com/DavideNale/distributed
```
2. Run it with make:
``` bash
make run
```

## API Usage

### Store a File

Stores a file with a specific key and broadcasts the update to all peers.

```go
func (s *FileServer) Store(key string, r io.Reader) error
```

**Parameters:**

- `key` (string): The unique identifier for the file.
- `r` (io.Reader): The file data to be stored.

**Example:**

```go
file, err := os.Open("example.txt")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

err = fileServer.Store("example_key", file)
if err != nil {
    log.Fatal(err)
}
```

### Get a File

Retrieves a file by its key, either locally or from a peer if not found locally.

```go
func (s *FileServer) Get(key string) (io.Reader, error)
```

**Parameters:**

- `key` (string): The unique identifier for the file.

**Example:**

```go
reader, err := fileServer.Get("example_key")
if err != nil {
    log.Fatal(err)
}

// Use the reader to read the file content
```

### Delete a File

Deletes a file by its key and broadcasts the deletion to all peers.

```go
func (s *FileServer) Delete(key string) error
```

**Parameters:**

- `key` (string): The unique identifier for the file.

**Example:**

```go
err := fileServer.Delete("example_key")
if err != nil {
    log.Fatal(err)
}
```

### Clear the File System

Deletes all files in the file system.

```go
func (s *FileServer) Clear() error
```

**Example:**

```go
err := fileServer.Clear()
if err != nil {
    log.Fatal(err)
}
```

## Internal Mechanisms

### Gossip Protocol

The gossip protocol is used to propagate changes across the network. Each file operation (store, get, delete) triggers a message to be broadcasted to all peers. The messages are structured as follows:

```go
type Message struct {
    Payload interface{}
}

type MessageStore struct {
    Key  string
    Size int64
}

type MessageGet struct {
    Key string
}

type MessageDelete struct {
    Key string
}
```

### TCP Streaming

TCP streams are used to transfer file data between peers. This ensures efficient and reliable data transfer. The `io.MultiWriter` is used to broadcast data to multiple peers simultaneously.

## Logging

The system uses a logger to provide detailed information about each operation. Log levels include `Info`, `Debug`, and `Warn`.

**Example:**

```go
s.Logger.Info("successfully stored file", "key", key, "size", size)
s.Logger.Debug("broadcasting to peers")
```