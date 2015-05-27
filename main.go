package main

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	repos  = []string{}
	args   = []string{"git", "status", "--short"}
	stdout = os.Stdout
	stderr = os.Stderr
	stdin  = os.Stdin
	color  = Blue
)

func main() {
	for _, gp := range filepath.SplitList(os.Getenv("GOPATH")) {
		for repo := range findRepoInGopath(gp) {
			blank, err := checkIfCmdReturnBlank(args, repo)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			if !blank {
				printColor(color, "$GOPATH"+repo[len(gp):])
				if err := run(args, color, repo); err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
			}
		}
	}
}
