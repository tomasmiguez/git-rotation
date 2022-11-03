package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func exPath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}

var dir = flag.String("dir", exPath(), "directory of the repository")

type Interval struct {
	From time.Time
	To   time.Time
}

func (interval Interval) duration() time.Duration {
	return interval.To.Sub(interval.From)
}

func fmtDuration(duration time.Duration) string {
	days := duration.Hours() / 24

	return fmt.Sprintf("%.f", days)
}

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

	intervals := make(map[string]Interval)
	err = cIter.ForEach(func(c *object.Commit) error {
		if entry, prs := intervals[c.Author.Email]; !prs {
			intervals[c.Author.Email] = Interval{From: c.Author.When, To: c.Author.When}
		} else {
			if c.Author.When.Before(entry.From) {
				entry.From = c.Author.When
			}
			if c.Author.When.After(entry.To) {
				entry.To = c.Author.When
			}
			intervals[c.Author.Email] = entry
		}

		return nil
	})
	for name, interval := range intervals {
		fmt.Println(name, ": ", fmtDuration(interval.duration()), " (", interval.From.Format("2006/01/02"), " ", interval.To.Format("2006/01/02"), ")")
	}
}
