/*
duct - add stdin and stdout to a code formatter

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
*/
package main
