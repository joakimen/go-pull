package main

import (
	"fmt"
	"io/ioutil"
	"kajito/git"
	"os"
	"path/filepath"
	"sync"
)

func main() {

	root, ok := os.LookupEnv("REPO")

	// test if $REPO is set
	if !ok {
		fmt.Println("Couldn't env var $REPO")
		return
	}

	dirs, err := ioutil.ReadDir(root)
	if err != nil {
		fmt.Println(err)
		return
	}

	var wg sync.WaitGroup
	var results []git.PullResult
	for _, dir := range dirs {

		absDir := filepath.Join(root, dir.Name())
		repo, err := git.New(absDir)

		// if err != nil then the path is not a valid git repo.
		if err != nil {
			continue
		}

		// update repo
		wg.Add(1)
		go func() {
			defer wg.Done()
			result := repo.Pull()

			if len(result.Commits) > 0 {
				results = append(results, result)
			}
		}()
	}

	wg.Wait()

	// Output results
	for _, res := range results {
		fmt.Println(res.Repo)   // print repo-name
		printSlice(res.Commits) // print unmerged commits
	}

}

func printSlice(s []string) {
	for _, e := range s {
		if e != "" {
			fmt.Println("- ", e)
		}
	}
}
