package main

import (
	"os"
	"path/filepath"
	"strings"
)

func findRepoInPath(target string) chan string {
	repos := make(chan string)
	go func() {
		filepath.Walk(target, func(path string, info os.FileInfo, err error) error {
			if info == nil {
				return err
			}
			if !info.IsDir() {
				return err
			}
			if strings.HasPrefix(info.Name(), ".") && target != path {
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
