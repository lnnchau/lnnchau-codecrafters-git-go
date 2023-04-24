package main

import (
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func getObjectFn(sha1Hash string) string {
	folder := sha1Hash[:2]
	file := sha1Hash[2:]
	objectFile := fmt.Sprintf(".git/objects/%s/%s", folder, file)
	return objectFile
}

func getStoreObject(objectHash string) string {
	f, err := os.Open(getObjectFn(objectHash))
	check(err)
	defer f.Close()

	r, err := zlib.NewReader(f)
	check(err)
	defer r.Close()

	bytesContents, err := io.ReadAll(r)
	check(err)

	store := string(bytesContents)
	return store
}

func parseHeader(header string) (string, int) {
	headerParts := strings.Split(header, " ")
	objectType := headerParts[0]
	objectSize, err := strconv.Atoi(headerParts[1])

	check(err)

	return objectType, objectSize
}

func WriteTree(folder string) (string, error) {
	sha1Hash := ""
	contentString := ""

	files, err := ioutil.ReadDir(folder)
	check(err)

	for _, file := range files {
		if file.Name() == ".git" {
			continue
		}

		if file.IsDir() {
			sha1Hash, err = WriteTree(filepath.Join(folder, file.Name()))
			contentString += fmt.Sprintf("40000 tree %s %s\n", sha1Hash, file.Name())
		} else {
			sha1Hash, err = WriteBlob(filepath.Join(folder, file.Name()))
			contentString += fmt.Sprintf("100644 blob %s %s\n", sha1Hash, file.Name())
		}
	}

	check(err)
	sha1Hash, err = writeHashFile(fmt.Sprintf("tree %d\u0000%s", len(contentString), contentString))
	return sha1Hash, err
}

func WriteBlob(path string) (string, error) {
	f, err := os.Open(path)
	check(err)
	defer f.Close()

	content, err := io.ReadAll(f)
	check(err)

	contentString := fmt.Sprintf("blob %d\u0000%s", len(content), content)

	// create sha1 hash
	// create objectFile in file system
	// write content to objectFile
	sha1Hash, err := writeHashFile(contentString)

	return sha1Hash, err
}

func writeHashFile(contentString string) (string, error) {
	sha1Hash := fmt.Sprintf("%x", sha1.Sum([]byte(contentString)))
	objectFile := getObjectFn(sha1Hash)

	// check if objectFile exists
	// if exists, return sha1Hash
	_, err := os.Stat(objectFile)
	if err == nil {
		return sha1Hash, nil
	}

	check(os.MkdirAll(filepath.Dir(objectFile), os.ModePerm))
	wF, err := os.Create(objectFile)
	check(err)
	defer wF.Close()

	w := zlib.NewWriter(wF)
	defer w.Close()

	_, err = w.Write([]byte(contentString))
	return sha1Hash, err
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
