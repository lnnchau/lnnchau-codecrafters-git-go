package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

type Tree struct {
	path       string
	fileMode   string
	objectType string
	objects    []GitObject
}

func NewTree(path string) *Tree {
	tree := &Tree{
		path:       path,
		fileMode:   "40000",
		objectType: "tree",
	}

	files, err := ioutil.ReadDir(tree.path)
	check(err)

	for _, file := range files {
		if file.Name() == ".git" {
			continue
		}

		tree.addObject(NewGitObject(filepath.Join(tree.path, file.Name())))
	}

	return tree
}

func (t *Tree) GetObjectType() string {
	return t.objectType
}

func (t *Tree) addObject(newObject GitObject) {
	t.objects = append(t.objects, newObject)
}

func (t *Tree) GetString() string {
	_, hash, err := t.WriteToFile()
	check(err)

	return fmt.Sprintf("%s %v\u0000%s", t.fileMode, filepath.Base(t.path), hash)
}

func (t *Tree) WriteToFile() (string, [20]byte, error) {
	sha1Hash := ""
	contentString := ""

	for _, obj := range t.objects {
		contentString += obj.GetString()
	}

	sha1Hash, sha1HashInBytes, err := writeHashFile(GetContentString(t.objectType, string(contentString)))
	return sha1Hash, sha1HashInBytes, err
}
