package main

import (
	"compress/zlib"
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
