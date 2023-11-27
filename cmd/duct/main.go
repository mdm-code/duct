package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/mdm-code/duct"
)

var (
	usage = `duct - add stdin and stdout to a code formatter

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

The program wraps a code formatter, which accepts file names as command
arguments instead of reading from standard input data stream, inside of a
standard Unix stdin to stdout filter-like data flow. The -stdout and -stderr
flags replace stdout and stderr of duct with stdout and stderr of the wrapped
command. It proves useful when the wrapped command reads code from files but
writes its output to stdout and/or stderr and not directly to provided files.
`

	attachStdout bool
	attachStderr bool

	cmdStdout io.Writer = duct.Discard
	cmdStderr io.Writer = duct.Discard
)

func main() {
	os.Exit(run())
}

func run() int {
	nonFlagArgs := parseArgs()
	if len(nonFlagArgs) < 1 {
		fmt.Fprintf(os.Stderr, "ERROR: missing command to wrap")
		return 1
	}
	f, err := os.CreateTemp("", duct.Pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to create a temporary file: %s", err)
		return 1
	}
	ductFDs, closer := duct.NewFDs(os.Stdin, os.Stdout, duct.Discard, f)
	defer os.Remove(f.Name())
	defer closer()
	cmdName, cmdArgs := nonFlagArgs[0], nonFlagArgs[1:]
	cmdArgs = append(cmdArgs, f.Name())
	if attachStdout {
		cmdStdout = os.Stdout
	}
	if attachStderr {
		cmdStderr = os.Stderr
	}
	cmd := duct.Cmd(cmdName, cmdStdout, cmdStderr, cmdArgs...)
	if attachStdout || attachStderr {
		err = duct.WrapWrite(cmd, ductFDs)
	} else {
		err = duct.Wrap(cmd, ductFDs)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to reformat the file: %s", err)
		return 1
	}
	return 0
}

func parseArgs() []string {
	w := flag.CommandLine.Output()
	flag.Usage = func() {
		fmt.Fprint(w, usage)
	}
	flag.BoolVar(&attachStdout, "stdout", attachStdout, "")
	flag.BoolVar(&attachStderr, "stderr", attachStderr, "")
	flag.Parse()
	return flag.Args()
}
