package main

import (
	"os"
	"path/filepath"
	"strings"
)

func listRepoInGopath() {
	var quit chan bool
	quit = make(chan bool)
	go func() {
		cwd := filepath.Join(gopath, "src")
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
			repos = append(repos, path)
			return nil
		})
		quit <- true
	}()
	if quit != nil {
		<-quit
	}
}
