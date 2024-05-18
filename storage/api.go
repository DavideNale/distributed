package storage

import (
	"errors"
	"fmt"
	"io"
	"os"
)

// Write writes data from the provided io.Reader to the store under the given key.
func (s *Store) Write(key string, r io.Reader) (int64, error) {
	return s.writeStream(key, r)
}

// Read reads data associated with the provided key from the store.
func (s *Store) Read(key string) (int64, io.Reader, error) {
	return s.readStream(key)
}

// Delete deletes the data associated with the provided key from the store.
func (s *Store) Delete(key string) error {
	path := s.PathTransformer(key)
	filePath := fmt.Sprintf("%s/%s", s.Root, path.FullPath())
	return os.RemoveAll(filePath)
}

// Exists checks if data associated with the provided key exists in the store.
func (s *Store) Exists(key string) bool {
	path := s.PathTransformer(key)
	filePath := fmt.Sprintf("%s/%s", s.Root, path.FullPath())
	_, err := os.Stat(filePath)
	return !errors.Is(err, os.ErrNotExist)
}

// Clear removes all data from the store, leaving it empty.
func (s *Store) Clear() error {
	return os.RemoveAll(s.Root)
}
