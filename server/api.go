package server

import (
	"bytes"
	"io"
)

func (s *FileServer) Store(key string, r io.Reader) error {
	fileBuffer := new(bytes.Buffer)
	reader := io.TeeReader(r, fileBuffer)

	size, err := s.store.Write(key, reader)
	_ = size
	if err != nil {
		return err
	}
	return nil
}
