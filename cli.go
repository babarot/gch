package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
	ExitCodeNotBlank
	ExitCodeRunError
)

type CLI struct {
	outStream, errStream io.Writer
}

type Config struct {
	Repos []string
}

var (
	repos  = []string{}
	git    = []string{"git", "status", "--short"}
	stdout = os.Stdout
	stderr = os.Stderr
	stdin  = os.Stdin
	color  = Blue
)

func (cli *CLI) Run(args []string) int {
	cpu := runtime.NumCPU()
	runtime.GOMAXPROCS(cpu)

	for _, gp := range filepath.SplitList(os.Getenv("GOPATH")) {
		for repo := range findRepoInGopath(gp) {
			blank, err := checkIfCmdReturnBlank(git, repo)
			if err != nil {
				fmt.Fprintf(cli.errStream, ColoredError(err.Error()))
				return ExitCodeNotBlank
			}
			if !blank {
				printColor(color, "$GOPATH"+repo[len(gp):])
				if err := run(git, color, repo); err != nil {
					fmt.Fprintf(cli.errStream, ColoredError(err.Error()))
					return ExitCodeRunError
				}
			}
		}
	}

	return ExitCodeOK
}
