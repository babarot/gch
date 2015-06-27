package main

import (
	"flag"
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
	ExitCodeErrorBlank
	ExitCodeErrorRunCommand
	ExitCodeErrorParseFlag
)

type CLI struct {
	outStream, errStream io.Writer
}

type Config struct {
	Repos []string
}

type Target struct {
	Path    string
	Channel chan string
}

var (
	command = []string{"git", "status", "--short"}
	stdout  = os.Stdout
	stderr  = os.Stderr
	stdin   = os.Stdin
	color   = Blue
	failed  = Red
)

func (c *CLI) Run(args []string) int {
	var list bool
	flags := flag.NewFlagSet("gch", flag.ContinueOnError)
	flags.SetOutput(c.errStream)
	flags.Usage = func() {
		fmt.Fprint(c.errStream, helpText)
	}
	flags.BoolVar(&list, "list", false, "")
	flags.BoolVar(&list, "l", false, "")

	if err := flags.Parse(args); err != nil {
		return ExitCodeErrorParseFlag
	}

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
			blank, err := checkIfCmdReturnBlank(command, repo)

			if err != nil {
				printColor(c.outStream, failed, err.Error())
				return ExitCodeErrorBlank
			}
			if !blank {
				printColor(c.outStream, color, repo)
				if list {
					continue
				}
				if err := runCommand(command, color, repo); err != nil {
					printColor(c.errStream, failed, err.Error())
					return ExitCodeErrorRunCommand
				}
			}
		}
	}

	return ExitCodeOK
}

var helpText = `Usage: gch [options]
  gch is a tool to run "git status" in every $GOPATHs recursively.

Options:
--list, -l    View only directory path in $GOPATHs
              without running git status.
`
