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

// Exit codes are int values that represent an exit code for a particular error.
// Sub-systems may check this unique error to determine the cause of an error
// without parsing the output or help text.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
	ExitCodeErrorBlank
	ExitCodeErrorRunCommand
	ExitCodeErrorParseFlag
	ExitCodeErrorGopathNotSet
)

// CLI is the command line object
type CLI struct {
	// outStream and errStream are the standard out and standard error streams to
	// write messages from the CLI.
	outStream, errStream io.Writer
}

type Target struct {
	Path    string
	Channel chan string
}

// Run invokes the CLI with the given arguments. The first argument is always
// the name of the application. This method slices accordingly.
func (cli *CLI) Run(args []string) int {
	var version, list bool

	flags := flag.NewFlagSet("gch", flag.ContinueOnError)
	flags.SetOutput(cli.errStream)
	flags.Usage = func() {
		fmt.Fprint(cli.errStream, helpText)
	}

	flags.BoolVar(&version, "version", false, "")
	flags.BoolVar(&list, "list", false, "")
	flags.BoolVar(&list, "l", false, "")

	// Parse all the flags
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeErrorParseFlag
	}

	// Version
	if version {
		fmt.Fprintf(cli.errStream, "%s v%s\n", Name, Version)
		return ExitCodeOK
	}

	var (
		command = []string{"git", "status", "--short"}
		color   = Blue
		failed  = Red
	)

	cpu := runtime.NumCPU()
	runtime.GOMAXPROCS(cpu)

	targets := []Target{}

	// Check if $GOPATH is set
	if os.Getenv("GOPATH") == "" {
		printColor(cli.outStream, failed, "$GOPATH not set")
		return ExitCodeErrorGopathNotSet
	}

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
				printColor(cli.outStream, failed, err.Error())
				return ExitCodeErrorBlank
			}
			if !blank {
				printColor(cli.outStream, color, repo)
				if list {
					continue
				}
				if err := runCommand(command, color, repo); err != nil {
					printColor(cli.errStream, failed, err.Error())
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
  --list, -l     View only directory path in $GOPATHs
                 without running git status.
  --version      Print the version of this application
`
