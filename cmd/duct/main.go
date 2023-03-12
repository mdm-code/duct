package main

import (
	"fmt"
	"io"
	"os"

	"github.com/mdm-code/duct"
)

func main() {
	f, err := duct.NewTempf()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create a temporary files: %s", err)
	}
	fds := duct.Fds{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: io.Discard,
		Tempf:  f,
	}
	cmd := os.Args[1]
	err = duct.Wrap(cmd, &fds)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to invoke %s: %s", cmd, err)
	}
}
