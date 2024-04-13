package storage

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
)

type PathInfo struct {
	FilePath string
	FileName string
}

func (p PathInfo) FullPath() string {
	return fmt.Sprintf("%s/%s", p.FilePath, p.FileName)
}

type PathTransformer func(string) PathInfo

var DefaultTransformer = func(key string) PathInfo {
	return PathInfo{
		FilePath: key,
		FileName: key,
	}
}

var HashTransformer = func(key string) PathInfo {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blocksize := 5
	slices := len(hashStr) / blocksize

	paths := make([]string, slices)
	for i := range slices {
		from, to := i*blocksize, (i+1)*blocksize
		paths[i] = hashStr[from:to]
	}

	return PathInfo{
		FilePath: strings.Join(paths, "/"),
		FileName: hashStr,
	}
}
