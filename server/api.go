package server

import (
	"bytes"
	"encoding/binary"
	"io"
	"time"
)

// Store stores a file with a specific key.
func (s *FileServer) Store(key string, r io.Reader) error {
	fileBuffer := new(bytes.Buffer)
	reader := io.TeeReader(r, fileBuffer)

	size, err := s.store.Write(key, reader)
	if err != nil {
		return err
	}

	s.Logger.Info("successfully stored file", "key", key, "size", size)
	s.Logger.Debug("broadcasting to peers")

	msg := Message{
		Payload: MessageStore{
			Key:  key,
			Size: size,
		},
	}

	if err := s.broadcast(&msg); err != nil {
		return err
	}

	time.Sleep(time.Microsecond * 5)

	peers := []io.Writer{}
	for _, peer := range s.peers {
		peers = append(peers, peer)
	}
	mw := io.MultiWriter(peers...)
	_, err = mw.Write(fileBuffer.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func (s *FileServer) Get(key string) (io.Reader, error) {
	if s.store.Exists(key) {
		size, r, err := s.store.Read(key)
		if err != nil {
			return nil, err
		}
		s.Logger.Info("successfully retrieved local file", "key", key, "size", size)
		return r, nil
	}
	s.Logger.Debug("file not present locally, fetching for peers", "key", key)

	msg := Message{
		Payload: MessageGet{
			Key: key,
		},
	}

	if err := s.broadcast(&msg); err != nil {
		return nil, err
	}

	for _, peer := range s.peers {
		var size int64
		binary.Read(peer, binary.LittleEndian, &size)
		n, err := s.store.Write(key, io.LimitReader(peer, size))
		if err != nil {
			return nil, err
		}
		s.Logger.Info("recieved from peer", "key", key, "size", n)
		peer.CloseStream()
		break
	}

	time.Sleep(time.Millisecond * 100)
	return nil, nil
}

// Delete deletes the file with the specified key, if it exists.
func (s *FileServer) Delete(key string) error {
	defer s.Logger.Info("deleted file", "key", key)

	msg := Message{
		Payload: MessageDelete{
			Key: key,
		},
	}

	if err := s.broadcast(&msg); err != nil {
		return err
	}
	return s.store.Delete(key)
}

// Clear deletes all inside the root of the file system.
func (s *FileServer) Clear() error {
	defer s.Logger.Warn("file system cleanup successful")
	return s.store.Clear()
}
