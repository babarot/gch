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

type Target struct {
	Path string
	Channel chan string
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

	targets := []Target{}

	// Add $GOPATH into target if exist
	for _, line := range filepath.SplitList(os.Getenv("GOPATH")) {
		goPath := filepath.Join(line, "src")
		existFlag := false
		for _, t := range targets {
			if t.Path == goPath {
				existFlag = true
			}
		}
		if existFlag {
			break
		}
		target := Target{goPath, findRepoInPath(goPath)}
		targets = append(targets, target)
	}

	// Add 'ghq root' into target if exist
	out, err := exec.Command("ghq", "root").Output()
	if err == nil {
		ghqPath := string(out)[:len(out)-1]
		existFlag := false
		for _, t := range targets {
			if t.Path == ghqPath {
				existFlag = true
			}
		}
		if !existFlag {
			target := Target{ghqPath, findRepoInPath(ghqPath)}
			targets = append(targets, target)
		}
	}

	for _, target := range targets {
		for repo := range target.Channel {
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
