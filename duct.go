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

// NameReadWriteSeekCloser specifies the interface for the temporary file. On
// top of a set of standard IO methods, it adds Name() used to retrieve the
// name of the file passed to the wrapped shell command.
type NameReadWriteSeekCloser interface {
	io.Reader
	io.Writer
	io.Seeker
	io.Closer
	Name() string
}

// FDs groups file descriptors used in the process of shell command wrapping.
type FDs struct {
	Stdin          io.ReadCloser
	Stdout, Stderr io.WriteCloser
	TempFile       NameReadWriteSeekCloser
}

type discard struct{}

// Close consecutively calls Close() on all file descriptors.
func (f *FDs) Close() error {
	for _, c := range []io.Closer{f.Stdin, f.Stdout, f.Stderr, f.TempFile} {
		if c == nil {
			return NilFDError
		}
		err := c.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// NewFDs groups file descriptors passed as function arguments in a single
// struct.
func NewFDs(stdin io.ReadCloser, stdout, stderr io.WriteCloser, tempFile NameReadWriteSeekCloser) *FDs {
	return &FDs{
		Stdin:    stdin,
		Stdout:   stdout,
		Stderr:   stderr,
		TempFile: tempFile,
	}
}

// Write on the discard struct always succeeds when it is invoked.
func (discard) Write(b []byte) (int, error) { return len(b), nil }

// Close on the discard struct never raises an error when it is invoked.
func (discard) Close() error { return nil }

//
func Cmd(name string, stdout, stderr io.Writer, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Stdout, cmd.Stderr = stdout, stderr
	return cmd
}

// Wrap invokes a code formatter by its name with source code read from
// standard input.
//
// Code to be formatted is being read from the fds.Stdin and written to
// fds.Stdout with fds.TempFile read/write functioning as an intermediate step
// necessitated by the design of the CLI interface of the formatter.
func Wrap(name string, fds *FDs) error {
	in := bufio.NewReader(fds.Stdin)
	_, err := in.WriteTo(fds.TempFile)
	if err != nil {
		return err
	}
	_, err = fds.TempFile.Seek(0, 0)
	if err != nil {
		return err
	}
	cmd := Cmd(name, fds.Stdout, fds.Stderr, fds.TempFile.Name())
	err = cmd.Run()
	if err != nil {
		return err
	}
	out := bufio.NewWriter(fds.Stdout)
	_, err = out.ReadFrom(fds.TempFile)
	if err != nil {
		return err
	}
	return nil
}
