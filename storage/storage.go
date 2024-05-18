package storage

import (
	"fmt"
	"io"
	"os"
)

type StoreOpts struct {
	Root            string // The root folder of the file system
	PathTransformer PathTransformer
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	path := s.PathTransformer(key)
	return os.Open(s.Root + "/" + path.FullPath())
}

func (s *Store) writeStream(key string, r io.Reader) (int64, error) {
	path := s.PathTransformer(key)
	folderPath := fmt.Sprintf("%s/%s", s.Root, path.FilePath)
	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		return 0, err
	}

	filePath := fmt.Sprintf("%s/%s", s.Root, path.FullPath())
	f, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}

	return io.Copy(f, r)
}
