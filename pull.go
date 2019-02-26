package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sync"

	"github.com/krystah/git"
	"github.com/pkg/errors"
)

func main() {

	fRoot := flag.String("root", "", "Base directory containing the repositories")
	flag.Parse()
	root := *fRoot

	// if root-argument was not supplied, return an error
	if root == "" {
		log.Fatalf("-root was not specified")
	}

	err := iterate(root)

	if err != nil {
		log.Fatalf("error while pulling repository %s: ", err)
	}
}

// Iterate updates all repositories contained in "root"
func iterate(root string) error {

	type result struct {
		repo    string
		commits []string
	}

	var results []result

	dirs, err := ioutil.ReadDir(root)
	if err != nil {
		return errors.Wrap(err, "read failed")
	}

	var wg sync.WaitGroup
	var commits []string
	for _, dir := range dirs {

		repo := filepath.Join(root, dir.Name())
		if !git.IsValidRepo(repo) {
			continue
		}
		// update repo
		wg.Add(1)
		go func(repo string) {
			defer wg.Done()
			commits, _ = git.Pull(repo)

			if len(commits) > 0 {
				results = append(results, result{dir.Name(), commits})
			}
		}(repo)
	}

	// wait for all goroutines to finish
	wg.Wait()

	// print changes
	for _, r := range results {
		fmt.Println(r.repo) // print repo-name
		for _, c := range r.commits {
			fmt.Printf("* %s\n", c)
		}
	}
	return nil
}
