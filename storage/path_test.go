package storage

import (
	"fmt"
	"testing"
)

func TestTransformer(t *testing.T) {
	key := "test_key"
	pathInfo := HashTransformer(key)
	fmt.Printf("%+v\n", pathInfo)
}
