package storage

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
)

// Write
// Read
// Delete
// Exists
// Clear

func (s *Store) Write(key string, r io.Reader) error {
	return s.writeStream(key, r)
}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)

	return buf, err
}

func (s *Store) Delete(key string) error {
	path := s.PathTransformer(key)
	return os.RemoveAll(s.Root + "/" + path.FullPath())
}

func (s *Store) Exists(key string) bool {
	path := s.PathTransformer(key)

	_, err := os.Stat(s.Root + "/" + path.FullPath())
	if errors.Is(err, fs.ErrNotExist) {
		return false
	}
	return true
}

func (s *Store) Clear() error {
	return os.RemoveAll(s.Root)
}
