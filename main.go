package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
  "github.com/go-git/go-git/v5/plumbing/object"
)

func exPath() (string) {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}

var dir = flag.String("dir", exPath(), "directory of the repository")

func main() {
  flag.Parse()

  repo, err := git.PlainOpen(*dir)
  if err != nil {
    panic(err)
  }

  cIter, err := repo.Log(&git.LogOptions{})
  if err != nil {
    panic(err)
  }

  var cCount int
  err = cIter.ForEach(func(c *object.Commit) error {
		cCount++

		return nil
	})

  fmt.Println(cCount)
}
