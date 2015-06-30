package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestRun__versionFlag(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}
	args := strings.Split("gch --version", " ")

	status := cli.Run(args)
	if status != ExitCodeOK {
		t.Errorf("expected %d to eq %d", status, ExitCodeOK)
	}

	expected := fmt.Sprintf("%s v%s", Name, Version)
	if !strings.Contains(errStream.String(), expected) {
		t.Errorf("expected %q to eq %q", errStream.String(), expected)
	}
}

func TestRun_parseError(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}
	args := strings.Split("gch --unknown-flag", " ")

	status := cli.Run(args)
	if status != ExitCodeErrorParseFlag {
		t.Errorf("expected %d to eq %d", status, ExitCodeErrorParseFlag)
	}

	expected := "flag provided but not defined: -unknown-flag"
	if !strings.Contains(errStream.String(), expected) {
		t.Fatalf("expected %q to contain %q", errStream.String(), expected)
	}
}
