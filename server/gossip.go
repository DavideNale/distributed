package server

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
)

type Message struct {
	Payload any
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

func init() {
	gob.Register(MessageStore{})
	gob.Register(MessageDelete{})
	gob.Register(MessageGet{})
}

func (s *FileServer) broadcast(msg *Message) error {
	buf := new(bytes.Buffer)

	if err := gob.NewEncoder(buf).Encode(msg); err != nil {
		return err
	}

	for _, peer := range s.peers {
		if err := peer.Send(buf.Bytes()); err != nil {
			return err
		}
	}
	return nil
}

func (s *FileServer) handleMessage(from string, msg *Message) error {
	switch v := msg.Payload.(type) {
	case MessageStore:
		return s.handleMessageStore(from, v)
	case MessageGet:
		return s.handleMessageGet(from, v)
	case MessageDelete:
		return s.handleMessageDelete(from, v)
	}
	return nil
}

func (s *FileServer) handleMessageStore(from string, msg MessageStore) error {
	peer, ok := s.peers[from]
	if !ok {
		s.Logger.Warn("peer not in peer list", "peer", peer.LocalAddr())
	}

	n, err := s.store.Write(msg.Key, io.LimitReader(peer, msg.Size))
	if err != nil {
		return nil
	}
	s.Logger.Info("recieved file from peer", "key", msg.Key, "size", n)
	peer.CloseStream()
	return nil
}

func (s *FileServer) handleMessageGet(from string, msg MessageGet) error {
	peer, ok := s.peers[from]
	if !ok {
		s.Logger.Warn("peer not in peer list", "peer", peer.LocalAddr())
	}
	if !s.store.Exists(msg.Key) {
		return fmt.Errorf("file not found locally")
	}

	size, r, err := s.store.Read(msg.Key)
	if err != nil {
		return err
	}

	if rc, ok := r.(io.ReadCloser); ok {
		defer rc.Close()
	}

	binary.Write(peer, binary.LittleEndian, size)
	n, err := io.Copy(peer, r)
	if err != nil {
		return nil
	}

	s.Logger.Debug("sent file over to peer", "key", msg.Key, "size", n)
	return nil
}

func (s *FileServer) handleMessageDelete(from string, msg MessageDelete) error {
	peer, ok := s.peers[from]
	if !ok {
		s.Logger.Warn("peer not in peer list", "peer", peer.LocalAddr())
	}
	return s.store.Delete(msg.Key)
}
