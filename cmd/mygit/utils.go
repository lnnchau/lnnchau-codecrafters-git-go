package main

import (
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
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

func writeHashFile(contentString string) (string, [20]byte, error) {
	sha1HashInBytes := sha1.Sum([]byte(contentString))
	sha1Hash := fmt.Sprintf("%x", sha1HashInBytes)
	objectFile := getObjectFn(sha1Hash)

	// check if objectFile exists
	// if exists, return sha1Hash
	_, err := os.Stat(objectFile)
	if err == nil {
		return sha1Hash, sha1HashInBytes, nil
	}

	check(os.MkdirAll(filepath.Dir(objectFile), os.ModePerm))
	wF, err := os.Create(objectFile)
	check(err)
	defer wF.Close()

	w := zlib.NewWriter(wF)
	defer w.Close()

	_, err = w.Write([]byte(contentString))
	return sha1Hash, sha1HashInBytes, err
}

func GetContentString(objectType string, content string) string {
	return fmt.Sprintf("%s %d\u0000%s", objectType, len(content), content)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
