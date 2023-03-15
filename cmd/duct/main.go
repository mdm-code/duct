package main

import (
	"fmt"
	"os"

	"github.com/mdm-code/duct"
)

var usage = `duct - add stdin and stdout to a code formatter

Duct wraps a code formatter inside of a stdin to stdout filter-like data flow.

Usage:
	duct [args...]

Options:
	-h, --help  show this help message and exit

Example:
	duct black -l 79 <<EOF
	from typing import (
		Protocol
	)
	class Sized(Protocol):
		def __len__(self) -> int: ...
	def print_size(s: Sized) -> None: len(s)
	class Queue:
		def __len__(self) -> int: return 10
	q = Queue(); print_size(q)
	EOF

Output:
	from typing import Protocol


	class Sized(Protocol):
		def __len__(self) -> int:
			...


	def print_size(s: Sized) -> None:
		len(s)


	class Queue:
		def __len__(self) -> int:
			return 10


	q = Queue()
	print_size(q)

The program wraps a code formatter, which accepts file names as commands
arguments instead of reading from standard input data stream, inside of a
standard Unix stdin to stdout filter-like data flow.
`

func main() {
	os.Exit(run())
}

func run() int {
	if len(os.Args) == 1 || (len(os.Args) > 1 && isHelp(os.Args[1])) {
		fmt.Fprintf(os.Stdout, usage)
		return 0
	}
	f, err := os.CreateTemp("", duct.Pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create a temporary file: %s", err)
		return 1
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
		return 1
	}
	return 0
}

func isHelp(arg string) bool {
	for _, flag := range []string{"-h", "-help", "--help"} {
		if arg == flag {
			return true
		}
	}
	return false
}
