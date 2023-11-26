package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/mdm-code/duct"
)

var usage = `duct - add stdin and stdout to a code formatter

Duct wraps a code formatter inside of a stdin to stdout filter-like data flow.

Usage:
	duct [OPTIONS] [args...]

Options:
	-h, -help, --help  show this help message and exit
	-stdout, --stdout  attach stdout of the wrapped command
	-stderr, --stderr  attach stderr of the wrapped command

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
	args := parseArgs()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "no command provided")
		return 1
	}
	f, err := os.CreateTemp("", duct.Pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create a temporary file: %s", err)
		return 1
	}
	fds, closer := duct.NewFDs(os.Stdin, os.Stdout, duct.Discard, f)
	defer os.Remove(f.Name())
	defer closer()
	cmdName, args := args[0], args[1:]
	args = append(args, f.Name())
	var o io.Writer = duct.Discard
	if stdout {
		o = os.Stdout
	}
	var e io.Writer = duct.Discard
	if stderr {
		e = os.Stderr
	}
	cmd := duct.Cmd(cmdName, o, e, args...)
	if stdout || stderr {
		err = duct.WrapXXX(cmd, fds)
	} else {
		err = duct.Wrap(cmd, fds)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to reformat the file: %s", err)
		return 1
	}
	return 0
}

var stdout bool
var stderr bool

// NOTE: consider using boolfunc that replaces discard with std in/out
// NOTE: the second occurrence will be added to unparsed args correctly,
// but it requires the 2nd one to be placed as the wrapped command arg like
// so:
// duct --stdout black --stdout=true
//
// In case the --stdout lands in front of the black, it will be evaluated
// as the duct native command argument.
func parseArgs() []string {
	w := flag.CommandLine.Output()
	flag.Usage = func() {
		fmt.Fprintf(w, usage)
	}
	flag.BoolVar(&stdout, "stdout", false, "")
	flag.BoolVar(&stderr, "stderr", false, "")
	flag.Parse()
	return flag.Args()
}
