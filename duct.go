/*
Package duct provides the internals for the duct command-line program wrapping
code formatters that do not read from standard input data stream, and instead
they take file names as command arguments. The package offers components that
allow such commands to be wrapped inside of a standard Unix STDIN to STDOUT
filter-like data flow.

The general idea is that input data read from standard input is written to an
intermediate temporary file. The name of the file gets passed to as one of the
positional arguments of the named program to be executed. The modified contents
of the file are then re-read and written out the standard output. This way the
wrapped program can be used as a regular Unix filter.
*/
package duct

import (
	"bufio"
	"errors"
	"io"
	"os/exec"
)

// Pattern defines the name pattern for the temporary file.
const Pattern = `duct-*`

// Discard is a WriteCloser that does nothing when either Write or Close
// methods are invoked. Ever call succeeds.
var Discard io.WriteCloser = discard{}

// NilFDError indicates that a file descriptor for read/write operation is nil.
var NilFDError error = errors.New("nil file descriptor")

// ReadWriteSeekCloser specifies the interface for the temporary file. On
// top of a set of standard IO methods, it adds Name() used to retrieve the
// name of the file passed to the wrapped shell command.
type ReadWriteSeekCloser interface {
	io.Reader
	io.Writer
	io.Seeker
	io.Closer
}

// Runner defines the interface for shell process to execute.
type Runner interface {
	Run() error
}

// FDs groups file descriptors used in the process of shell command wrapping.
type FDs struct {
	Stdin          io.ReadCloser
	Stdout, Stderr io.WriteCloser
	TempFile       ReadWriteSeekCloser
}

// NewFDs groups file descriptors passed as function arguments in a single
// struct.
//
// The closer method returned alongside the struct should be deferred to ensure
// that all files are closed upon the termination of the program.
func NewFDs(stdin io.ReadCloser, stdout, stderr io.WriteCloser, tempFile ReadWriteSeekCloser) (*FDs, func() error) {
	fds := &FDs{
		Stdin:    stdin,
		Stdout:   stdout,
		Stderr:   stderr,
		TempFile: tempFile,
	}
	return fds, (*fds).Close
}

// Close consecutively calls Close() on all file descriptors.
func (f *FDs) Close() error {
	for _, c := range []io.Closer{f.Stdin, f.Stdout, f.Stderr, f.TempFile} {
		if c == nil {
			return NilFDError
		}
		c.Close()
	}
	return nil
}

type discard struct{}

// Write on the discard struct always succeeds when it is invoked.
func (discard) Write(b []byte) (int, error) { return len(b), nil }

// Close on the discard struct never raises an error when it is invoked.
func (discard) Close() error { return nil }

// Cmd returns the Cmd stuct to execute a given named program with given
// arguments and file descriptors attached.
func Cmd(name string, stdout, stderr io.Writer, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Stdout, cmd.Stderr = stdout, stderr
	return cmd
}

// Wrap executes a given named formatter program cmd.
//
// Code to be formatted is being read from the fds.Stdin and written to
// fds.Stdout with fds.TempFile read/write functioning as an intermediate step
// necessitated by the design of the CLI interface of the formatter.
func Wrap(cmd Runner, fds *FDs) error {
	in := bufio.NewReader(fds.Stdin)
	_, err := in.WriteTo(fds.TempFile)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	_, err = fds.TempFile.Seek(0, 0)
	if err != nil {
		return err
	}
	err = cmd.Run()
	if err != nil {
		return err
	}
	out := bufio.NewWriter(fds.Stdout)
	_, err = out.ReadFrom(fds.TempFile)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	return nil
}
