package storage

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
)

// PathInfo contains the folder and file name.
type PathInfo struct {
	FilePath string // Directory path
	FileName string // File name
}

// FullPath returns the full path by concatenating the FilePath and FileName fields
// of the PathInfo struct, separated by a forward slash.
func (p PathInfo) FullPath() string {
	return fmt.Sprintf("%s/%s", p.FilePath, p.FileName)
}

// FullPath returns the concatenation of FilePath and FileName
type PathTransformer func(string) PathInfo

// DefaultTransformer returns a PathInfo where both FilePath and FileName are set to the input key.
// This is a simple identity transformer.
var DefaultTransformer = func(key string) PathInfo {
	return PathInfo{
		FilePath: key,
		FileName: key,
	}
}

// HashTransformer returns a PathInfo with a SHA-1 hashed directory structure.
// The FileName is set to the full hash string.
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
