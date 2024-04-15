package storage

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestTransformer(t *testing.T) {
	key := "test_key"
	pathname := HashTransformer(key)
	fmt.Println(pathname)
}

func TestStore(t *testing.T) {
	s := NewStore().WithTransformer(HashTransformer)

	key := "test_key"
	file := []byte("test_content")

	if err := s.writeStream(key, bytes.NewReader(file)); err != nil {
		t.Error(err)
	}

	r, err := s.Read(key)
	if err != nil {
		t.Error(err)
	}

	b, _ := io.ReadAll(r)
	fmt.Println(string(b))

	s.Clear()
}
