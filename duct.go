package duct

import (
	"bufio"
	"errors"
	"io"
	"os/exec"
)

// Pattern
const Pattern = `duct-*`

// Discard
var Discard io.WriteCloser = discard{}

// NilFDError
var NilFDError error = errors.New("")

// NameReadWriteSeekCloser
type NameReadWriteSeekCloser interface {
	io.Reader
	io.Writer
	io.Seeker
	io.Closer
	Name() string
}

// FDs
type FDs struct {
	Stdin          io.ReadCloser
	Stdout, Stderr io.WriteCloser
	TempFile       NameReadWriteSeekCloser
}

type discard struct{}

// Close
func (f *FDs) Close() error {
	for _, c := range []io.Closer{f.Stdin, f.Stdout, f.Stderr, f.TempFile} {
		err := c.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// NewFDs
func NewFDs(stdin io.ReadCloser, stdout, stderr io.WriteCloser, tempFile NameReadWriteSeekCloser) *FDs {
	return &FDs{
		Stdin:    stdin,
		Stdout:   stdout,
		Stderr:   stderr,
		TempFile: tempFile,
	}
}

// Write
func (discard) Write(b []byte) (int, error) { return len(b), nil }

// Close
func (discard) Close() error { return nil }

//
func Cmd(name string, stdout, stderr io.Writer, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Stdout, cmd.Stderr = stdout, stderr
	return cmd
}

// Wrap
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
