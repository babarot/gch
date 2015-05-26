package main

import (
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"os"
	"path/filepath"
)

type Color int

const (
	None Color = Color(ct.None)
	Red  Color = Color(ct.Red)
	Blue Color = Color(ct.Blue)
)

var (
	repos  = []string{}
	args   = []string{"git", "status", "-s"}
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
				if err := run(args, Blue, repo); err != nil {
				}
			}
		}
	}
}

func printColor(c Color, text string) {
	ct.ChangeColor(ct.Color(c), true, ct.None, false)
	fmt.Println(text)
	ct.ResetColor()
}
