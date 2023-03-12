package main

import (
	"fmt"
	"os"

	"github.com/mdm-code/duct"
)

func main() {
	f, err := os.CreateTemp("", duct.Pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create a temporary file: %s", err)
		os.Exit(1)
	}
	fds := duct.NewFDs(os.Stdin, os.Stdout, duct.Discard, f)
	cmd := os.Args[1]
	err = duct.Wrap(cmd, fds)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}
}
