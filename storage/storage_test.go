package storage

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestStore(t *testing.T) {
	storeOpts := StoreOpts{
		Root:            "data",
		PathTransformer: HashTransformer,
	}
	s := NewStore(storeOpts)

	key := "test_key"
	file := []byte("test_content")

	if _, err := s.writeStream(key, bytes.NewReader(file)); err != nil {
		t.Error(err)
	}

	_, r, err := s.Read(key)
	if err != nil {
		t.Error(err)
	}

	b, _ := io.ReadAll(r)
	fmt.Println(string(b))

	s.Clear()
}
