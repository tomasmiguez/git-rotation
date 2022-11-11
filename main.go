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

type interval struct {
	From time.Time
	To   time.Time
}

func (i interval) duration() int {
	return int(i.To.Sub(i.From).Hours() / 24)
}

func formatDate(d time.Time) string {
	return d.Format("2006/01/02")
}

func getCommitIter(dir string) (object.CommitIter, error) {
	repo, err := git.PlainOpen(dir)
	if err != nil {
		return nil, err
	}

	cIter, err := repo.Log(&git.LogOptions{})
	if err != nil {
		return nil, err
	}

	return cIter, nil
}

type intervalMap map[string]interval

func (intervals *intervalMap) processDir(dir string) (error) {
	cIter, err := getCommitIter(dir)
	if err != nil {
		return err
	}

	cIter.ForEach(func(c *object.Commit) error {
		if entry, prs := (*intervals)[c.Author.Email]; !prs {
			(*intervals)[c.Author.Email] = interval{From: c.Author.When, To: c.Author.When}
		} else {
			if c.Author.When.Before(entry.From) {
				entry.From = c.Author.When
			}
			if c.Author.When.After(entry.To) {
				entry.To = c.Author.When
			}
			(*intervals)[c.Author.Email] = entry
		}

		return nil
	})

	return nil
}

var (
	upTo = flag.String("upTo", "", "last termination date, formatted as 08-11-2022")
	min  = flag.Int("min", 0, "minimum numbers of days worked")
)

func main() {
	flag.Parse()

	var upToDate time.Time
	if *upTo != "" {
		var err error
		upToDate, err = time.Parse("02-01-2006", *upTo)
		if err != nil {
			panic(err)
		}
	}

	dirs := flag.Args()

	intervals := make(intervalMap)
	for _, dir := range dirs {
		err := intervals.processDir(dir)
		if err != nil {
			panic(err)
		}
	}

	var total int
	var acum int
	for name, i := range intervals {
		if ((*min < i.duration()) && (*upTo == "" || interval{From: upToDate, To: i.To}.duration() < 0)) {
			acum++
			total += i.duration()

			fmt.Println(name, ": ",	fmt.Sprintf("%.d", i.duration()), " (", formatDate(i.From), " ", formatDate(i.To), ")")
		}
	}

	fmt.Println()
	fmt.Println(total/acum)
}
