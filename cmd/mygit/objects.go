package main

import "os"

type GitObject interface {
	GetObjectType() string

	WriteToFile() (string, [20]byte, error)
	GetString() string
}

func NewGitObject(path string) GitObject {
	fi, err := os.Stat(path)
	check(err)

	if fi.IsDir() {
		return NewTree(path)
	}

	return NewBlob(path)
}
