package main

import (
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"os"
	"os/exec"
)

type Color int

const (
	None Color = Color(ct.None)
	Red  Color = Color(ct.Red)
	Blue Color = Color(ct.Blue)
)

func printColor(c Color, text string) {
	ct.ChangeColor(ct.Color(c), true, ct.None, false)
	fmt.Println(text)
	ct.ResetColor()
}

func run(args []string, c Color, path string) error {
	if len(args) == 0 {
		return fmt.Errorf("%s: invalid arguments", args)
	}
	if err := os.Chdir(path); err != nil {
		return err
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Stdin = stdin
	ct.ChangeColor(ct.Color(c), true, ct.None, false)
	err := cmd.Run()
	ct.ResetColor()
	return err
}

func checkIfCmdReturnBlank(args []string, path string) (bool, error) {
	if len(args) == 0 {
		return false, fmt.Errorf("%s: invalid arguments", args)
	}
	if err := os.Chdir(path); err != nil {
		return false, fmt.Errorf("%s: cannot change directory", path)
	}
	out, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		return false, fmt.Errorf("%s: returns failed", args)
	}
	return len(out) == 0, nil
}
