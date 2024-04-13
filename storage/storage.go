package storage

import (
	"io"
	"os"
)

type Store struct {
	Root            string
	PathTransformer PathTransformer
}

func NewStore() *Store {
	return &Store{
		Root:            "data",
		PathTransformer: DefaultTransformer,
	}
}

func (s *Store) WithRootPath(root string) *Store {
	s.Root = root
	return s
}

func (s *Store) WithTransformer(pt PathTransformer) *Store {
	s.PathTransformer = pt
	return s
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	path := s.PathTransformer(key)
	return os.Open(path.FullPath())
}

func (s *Store) writeStream(key string, r io.Reader) error {
	path := s.PathTransformer(key)

	if err := os.MkdirAll(path.FilePath, os.ModePerm); err != nil {
		return err
	}

	full := path.FullPath()

	f, err := os.Create(full)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, r)
	if err != nil {
		return err
	}

	return nil
}
