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
	fds, closer := duct.NewFDs(os.Stdin, os.Stdout, duct.Discard, f)
	defer closer()
	args := []string{}
	if len(os.Args) > 1 {
		args = append(args, os.Args[2:]...)
	}
	args = append(args, f.Name())
	cmd := duct.Cmd(os.Args[1], os.Stdout, duct.Discard, args...)
	err = duct.Wrap(cmd, fds)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to reformat the file: %s", err)
		os.Exit(1)
	}
}
