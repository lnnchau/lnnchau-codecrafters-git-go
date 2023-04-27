package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	SHA_1_SIZE_IN_BYTES = 20
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

		store := getStoreObject(objectHash)
		storeParts := strings.Split(store, "\u0000")

		header := storeParts[0]
		content := storeParts[1]

		if objectType == "-p" {
			fmt.Print(content)
		} else if objectType == "-t" {
			fmt.Print(header)
		}

	case "ls-tree":
		if len(os.Args) != 4 {
			fmt.Fprintf(os.Stderr, "usage: mygit ls-tree --name-only <hash>\n")
			os.Exit(1)
		}

		treeHash := os.Args[3]
		store := getStoreObject(treeHash)

		// the structure of the contents stored in the tree object is as follows:
		// tree <size>\u0000<content>
		// where <content> is a no-space concatenation of: <size> <name>\u0000<sha-1 hash>
		// for example: 100644 README.md\u0000<sha-1 hash>100644 README.md\u0000<sha-1 hash>

		// the following code splits the store into its parts
		storeParts := strings.Split(store, "\u0000")

		header, _ := parseHeader(storeParts[0])
		if header != "tree" {
			fmt.Fprintf(os.Stderr, "Error: %s is not a tree object\n", treeHash)
			os.Exit(1)
		}

		contentParts := storeParts[1 : len(storeParts)-1] // the last one would contains <sha-1 hash> only

		displayContent := ""
		for i, contentPart := range contentParts {
			// the first item isn't appended by the sha-1
			var sizeWithPath string
			if i == 0 {
				sizeWithPath = contentPart
			} else {
				sizeWithPath = contentPart[SHA_1_SIZE_IN_BYTES:]
			}

			sizeWithPathParts := strings.Split(sizeWithPath, " ")
			path := sizeWithPathParts[1]

			displayContent += fmt.Sprintln(path)
		}

		fmt.Print(displayContent)
	case "hash-object":
		if len(os.Args) != 4 {
			fmt.Fprintf(os.Stderr, "usage: mygit hash-object -w <file>\n")
			os.Exit(1)
		}

		fileName := os.Args[3]

		fileNameBlob := NewBlob(fileName)

		sha1Hash, _, err := fileNameBlob.WriteToFile()
		check(err)
		fmt.Println(sha1Hash)

	case "write-tree":
		wd, _ := os.Getwd()

		wdTree := NewTree(wd)
		sha1Hash, _, err := wdTree.WriteToFile()
		check(err)

		fmt.Println(sha1Hash)
	case "commit-tree":
		treeHash := os.Args[2]
		parentHash := os.Args[4]
		message := os.Args[6]

		commit := NewCommit(treeHash, parentHash, message)
		sha1Hash, _, err := commit.WriteToFile()
		check(err)

		fmt.Println(sha1Hash)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
