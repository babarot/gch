package main

import (
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"os"
	"strings"
)

type Color int

const (
	None Color = Color(ct.None)
	Red  Color = Color(ct.Red)
	Blue Color = Color(ct.Blue)
)

var (
	gopath = os.Getenv("GOPATH")
	repos  = []string{}
	args   = []string{"git", "status", "-s"}
	stdout = os.Stdout
	stderr = os.Stderr
	stdin  = os.Stdin
	color  = Blue
)

func main() {
	listRepoInGopath()

	for _, path := range repos {
		blank, err := checkIfCmdReturnBlank(args, path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if !blank {
			printColor(color, strings.Replace(path, gopath, "$GOPATH", 1))
			if err := run(args, Blue, path); err != nil {
			}
		}
	}
}

func printColor(c Color, text string) {
	ct.ChangeColor(ct.Color(c), true, ct.None, false)
	fmt.Println(text)
	ct.ResetColor()
}
