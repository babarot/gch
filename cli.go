package main

import (
	"flag"
	"fmt"
	"io"
	"os"
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

	for _, gp := range filepath.SplitList(os.Getenv("GOPATH")) {
		for repo := range findRepoInGopath(gp) {
			blank, err := checkIfCmdReturnBlank(command, repo)
			if err != nil {
				printColor(c.outStream, failed, err.Error())
				return ExitCodeErrorBlank
			}
			if !blank {
				if list {
					fmt.Fprintln(c.outStream, repo)
				} else {
					printColor(c.outStream, color, "$GOPATH"+repo[len(gp):])
					if err := runCommand(command, color, repo); err != nil {
						printColor(c.outStream, failed, err.Error())
						return ExitCodeErrorRunCommand
					}
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
