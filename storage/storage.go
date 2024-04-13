package storage

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
	"strings"
)

type PathTransformer func(string) string

var DefaultTransformer = func(key string) string { return key }

var HashTransformer = func(key string) string {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blocksize := 5
	slices := 20 / blocksize

	paths := make([]string, slices)
	for i := range slices {
		from, to := i*blocksize, (i+1)*blocksize
		paths[i] = hashStr[from:to]
	}

	return strings.Join(paths, "/")
}

type Store struct {
	PathTransformer PathTransformer
}

func NewStore() *Store {
	return &Store{
		PathTransformer: DefaultTransformer,
	}
}

func (s *Store) WithTransformer(pt PathTransformer) *Store {
	s.PathTransformer = pt
	return s
}

func (s *Store) writeStream(key string, r io.Reader) error {
	path := s.PathTransformer(key)
	_ = path

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}
	filename := "testname"
	complete := path + "/" + filename

	f, err := os.Create(complete)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, r)
	if err != nil {
		return err
	}
	return nil

}
