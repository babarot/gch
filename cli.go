package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
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

	goPath := filepath.Join(os.Getenv("GOPATH"), "src")
	target := []chan string{
		findRepoInPath(goPath),
	}

	// Add 'ghq root' into target if exist
	out, err := exec.Command("ghq", "root").Output()
	if err == nil {
		ghqPath := string(out)[:len(out)-1]
		target = append(target, findRepoInPath(ghqPath))
	}

	for _, channel := range target {
		for repo := range channel {
			blank, err := checkIfCmdReturnBlank(git, repo)

			if err != nil {
				fmt.Fprintf(cli.errStream, ColoredError(err.Error()))
				return ExitCodeNotBlank
			}
			if !blank {
				printColor(color, repo)
				if err := run(git, color, repo); err != nil {
					fmt.Fprintf(cli.errStream, ColoredError(err.Error()))
					return ExitCodeRunError
				}
			}
		}
	}

	return ExitCodeOK
}
