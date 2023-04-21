package main

import (
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"strings"
)

// Usage: your_git.sh <command> <arg1> <arg2> ...
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	switch command := os.Args[1]; command {
	case "init":
		for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
			}
		}

		headFileContents := []byte("ref: refs/heads/master\n")
		if err := os.WriteFile(".git/HEAD", headFileContents, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
		}

		fmt.Println("Initialized git directory")

	case "cat-file":
		if len(os.Args) != 4 {
			fmt.Fprintf(os.Stderr, "usage: mygit cat-file <type> <hash>\n")
			os.Exit(1)
		}

		objectType := os.Args[2]
		objectHash := os.Args[3]

		folder := objectHash[:2]
		fileName := objectHash[2:]

		f, err := os.Open(fmt.Sprintf(".git/objects/%s/%s", folder, fileName))
		check(err)
		defer f.Close()

		r, err := zlib.NewReader(f)
		check(err)
		defer r.Close()

		bytesContents, err := io.ReadAll(r)
		check(err)

		store := string(bytesContents)
		storeParts := strings.Split(store, "\u0000")

		header := storeParts[0]
		content := storeParts[1]

		if objectType == "-p" {
			fmt.Print(content)
		} else if objectType == "-t" {
			fmt.Print(header)
		}

	case "hash-object":
		if len(os.Args) != 4 {
			fmt.Fprintf(os.Stderr, "usage: mygit hash-object <file>\n")
			os.Exit(1)
		}

		fileName := os.Args[3]

		f, err := os.Open(fileName)
		check(err)
		defer f.Close()

		content, err := io.ReadAll(f)
		check(err)

		contentString := fmt.Sprintf("blob %d\u0000%s", len(content), content)

		// create sha1 hash
		sha1Hash := fmt.Sprintf("%x", sha1.Sum([]byte(contentString)))
		// print sha1 hash to stdout
		fmt.Print(sha1Hash)

		folder := sha1Hash[:2]
		file := sha1Hash[2:]
		objectFile := fmt.Sprintf(".git/objects/%s/%s", folder, file)

		// create objectFile in file system
		check(
			os.MkdirAll(
				fmt.Sprintf(".git/objects/%s", folder),
				os.ModePerm))
		wF, err := os.Create(objectFile)
		check(err)
		defer wF.Close()

		// write content to objectFile
		w := zlib.NewWriter(wF)
		defer w.Close()

		_, err = w.Write([]byte(contentString))
		check(err)

	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
