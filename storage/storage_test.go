package storage

import (
	"bytes"
	"fmt"
	"testing"
)

func TestTransformer(t *testing.T) {
	key := "test_key"
	pathname := HashTransformer(key)
	fmt.Println(pathname)
}

func TestStore(t *testing.T) {
	s := NewStore()

	file := bytes.NewReader([]byte("test content"))
	if err := s.writeStream("test_folder", file); err != nil {
		t.Error(err)
	}
}
