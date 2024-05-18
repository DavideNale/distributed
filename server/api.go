package server

import (
	"bytes"
	"io"
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
	return nil, nil
}

// Delete deletes the file with the specified key, if it exists.
func (s *FileServer) Delete(key string) error {
	defer s.Logger.Info("deleted file", "key", key)
	return s.store.Delete(key)
}

// Clear deletes all inside the root of the file system.
func (s *FileServer) Clear() error {
	defer s.Logger.Warn("file system cleanup successful")
	return s.store.Clear()
}
