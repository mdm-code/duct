<h1 align="center">
  <div>
    <img src="https://raw.githubusercontent.com/mdm-code/mdm-code.github.io/main/duct_logo.png" alt="logo"/>
  </div>
</h1>

<h4 align="center">Wrap a code formatter inside of a STDIN to STDOUT filter-like data flow</h4>

<div align="center">
<p>
    <a href="https://github.com/mdm-code/duct/actions?query=workflow%3ACI">
        <img alt="Build status" src="https://github.com/mdm-code/duct/workflows/CI/badge.svg">
    </a>
    <a href="https://app.codecov.io/gh/mdm-code/duct">
        <img alt="Code coverage" src="https://codecov.io/gh/mdm-code/duct/branch/main/graphs/badge.svg?branch=main">
    </a>
    <a href="https://opensource.org/licenses/MIT" rel="nofollow">
        <img alt="MIT license" src="https://img.shields.io/github/license/mdm-code/duct">
    </a>
    <a href="https://goreportcard.com/report/github.com/mdm-code/duct">
        <img alt="Go report card" src="https://goreportcard.com/badge/github.com/mdm-code/duct">
    </a>
    <a href="https://pkg.go.dev/github.com/mdm-code/duct">
        <img alt="Go package docs" src="https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white">
    </a>
</p>
</div>

The `duct` program allows to wrap code formatters inside of a STDIN to STDOUT
filter-like data flow. It wraps a code formatter, which accepts file names as
commands arguments instead of reading from standard input data stream, inside
of a standard Unix STDIN to STDOUT filter-like data flow. Consult the [package
documentation](https://pkg.go.dev/github.com/mdm-code/duct) or see
[Usage](#usage) to see how it works.


## Installation

Install the package to use the command-line `duct` to wrap code formatters
inside of a STDIN to STDOUT filter-like data flow.

```sh
go install github.com/mdm-code/duct@latest
```

Although I don't really see the reason why one might want to do it, use the
following command to add the package to an existing project.

```sh
go get github.com/mdm-code/duct
```


## Usage

Type `duct -h` to get information and examples on how to use `duct` and get
some ideas on how to use it in your workflow.

This very basic example show how to wrap `black`, a code formatter for Python,
with `duct` to use it as if it was a regular Unix data filter. Here is a snippet:

```sh
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
```

The example uses heredoc to direct code to standard input of `duct` that is
going to be formatter with `black` with the max line length set to 79
characters. The output is going to be written to STDOUT accordingly. This lets
you use `black` in `vim` as if it was a regular filter command, which makes
life much easier for a regular Python dev.


## Development

Consult [Makefile](Makefile) to see how to format, examine code with `go vet`,
run unit test, run code linter with `golint` in order to get test coverage and
check if the package builds all right.

Remember to install `golint` before you try to run tests and test the build:

```sh
go install golang.org/x/lint/golint@latest
```


## License

Copyright (c) 2023 Michał Adamczyk.

This project is licensed under the [MIT license](https://opensource.org/licenses/MIT).
See [LICENSE](LICENSE) for more details.
