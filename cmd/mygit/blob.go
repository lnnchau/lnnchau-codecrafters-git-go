package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Blob struct {
	fileMode   string
	objectType string
	path       string
}

func NewBlob(path string) *Blob {
	return &Blob{
		path:       path,
		fileMode:   "100644",
		objectType: "blob",
	}
}

func (b *Blob) GetObjectType() string {
	return b.objectType
}

func (b *Blob) GetString() string {
	_, filename := filepath.Split(b.path)
	_, hash, err := b.WriteToFile()
	check(err)

	return fmt.Sprintf("%s %v\u0000%s", b.fileMode, filename, hash)
}

func (b *Blob) WriteToFile() (string, [20]byte, error) {
	f, err := os.Open(b.path)
	check(err)
	defer f.Close()

	content, err := io.ReadAll(f)
	check(err)

	sha1Hash, sha1HashInBytes, err := writeHashFile(GetContentString(b.objectType, string(content)))

	return sha1Hash, sha1HashInBytes, err
}
