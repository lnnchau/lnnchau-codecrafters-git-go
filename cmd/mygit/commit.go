package main

import (
	"fmt"
	"time"
)

type Commit struct {
	treeSHA    string
	parentSHA  string
	objectType string
	message    string

	// author
	name  string
	email string
}

func NewCommit(treeSHA string, parentSHA string, message string) *Commit {
	return &Commit{
		treeSHA:    treeSHA,
		parentSHA:  parentSHA,
		objectType: "commit",
		message:    message,

		name:  "John Doe",
		email: "john.doe@example.com",
	}
}

func (c *Commit) GetObjectType() string {
	return c.objectType
}

func (c *Commit) GetString() string {
	panic("Not implemented")
}

func (c *Commit) WriteToFile() (string, [20]byte, error) {
	sha1Hash := ""
	contentString := fmt.Sprintf("tree %s\n", c.treeSHA)
	if c.parentSHA != "" {
		contentString += fmt.Sprintf("parent %s\n", c.parentSHA)
	}

	now := time.Now()
	timeStamp := fmt.Sprintf("%d %s", now.Unix(), now.Format("-0700"))
	authorString := fmt.Sprintf("author %s <%s> %s\n", c.name, c.email, timeStamp)
	committerString := fmt.Sprintf("committer %s <%s> %s\n", c.name, c.email, timeStamp)

	contentString += authorString + committerString + "\n" + c.message + "\n"

	sha1Hash, sha1HashInBytes, err := writeHashFile(GetContentString(c.objectType, string(contentString)))
	return sha1Hash, sha1HashInBytes, err
}
