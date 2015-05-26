package main

import (
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"os"
	"path/filepath"
	//"sync"
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
	//wg     sync.WaitGroup
)

func main() {
	for _, gp := range filepath.SplitList(os.Getenv("GOPATH")) {
		//wg.Add(1)
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
	//wg.Wait()
}

func printColor(c Color, text string) {
	ct.ChangeColor(ct.Color(c), true, ct.None, false)
	fmt.Println(text)
	ct.ResetColor()
}
