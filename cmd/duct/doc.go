/*
duct - add stdin and stdout to a code formatter

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
command. It's useful when the wrapped command reads code from files but writes
its output to stdout and/or stderr instead of writing directly to files.
*/
package main
