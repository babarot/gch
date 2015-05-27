package main

import (
	"os"
	"path/filepath"
	"strings"
)

func findRepoInGopath(gp string) chan string {
	repos := make(chan string)
	go func() {
		cwd := filepath.Join(gp, "src")
		filepath.Walk(cwd, func(path string, info os.FileInfo, err error) error {
			if info == nil {
				return err
			}
			if !info.IsDir() {
				return err
			}
			if strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			path = filepath.ToSlash(path)
			isgit := filepath.Join(path, ".git")
			if _, err := os.Stat(isgit); err != nil {
				return nil
			}
			repos <- path
			return filepath.SkipDir
		})
		close(repos)
	}()
	return repos
}
