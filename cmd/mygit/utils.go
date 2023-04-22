package main

import (
	"compress/zlib"
	"fmt"
	"io"
	"os"
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

func check(e error) {
	if e != nil {
		panic(e)
	}
}
