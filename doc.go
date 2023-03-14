/*
Package duct provides the internals for the duct command-line program wrapping
code formatters that do not read from standard input data stream, and instead
they take file names as command arguments. The package offers components that
allow such commands to be wrapped inside of a standard Unix stdin to stdout
filter-like data flow.

The general idea is that input data read from standard input is written to an
intermediate temporary file. The name of the file gets passed to as one of the
positional arguments of the named program to be executed. The modified contents
of the file are then re-read and written out the standard output. This way the
wrapped program can be used as a regular Unix filter.
*/
package duct
