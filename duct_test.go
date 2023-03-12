package duct

import (
	"errors"
	"fmt"
	"io"
	"testing"
)

type (
	MockedCmd struct{}

	MockedCmdFail struct{}

	MockedReader struct{}

	MockedReaderFail struct{}

	MockedWriter struct{}

	MockedWriterFail struct{}

	MockedCloser struct{}

	MockedSeeker struct{}

	MockedSeekerFail struct{}

	MockedReadCloser struct {
		MockedReader
		MockedCloser
	}

	MockedReadCloserFail struct {
		MockedReaderFail
		MockedCloser
	}

	MockedWriteCloser struct {
		MockedWriter
		MockedCloser
	}

	MockedWriteCloserFail struct {
		MockedWriterFail
		MockedCloser
	}

	MockedReadWriteSeekCloser struct {
		MockedReader
		MockedWriter
		MockedSeeker
		MockedCloser
	}

	MockedReadWriteSeekCloserFail struct {
		MockedReader
		MockedWriter
		MockedSeekerFail
		MockedCloser
	}

	MockedReadWriteSeekCloserFail2 struct {
		MockedReader
		MockedWriterFail
		MockedSeeker
		MockedCloser
	}
)

func (MockedCmd) Run() error { return nil }

func (MockedCmdFail) Run() error { return fmt.Errorf("errored") }

func (MockedReader) Read(b []byte) (int, error) { return len(b), io.EOF }

func (MockedReaderFail) Read(b []byte) (int, error) { return len(b), fmt.Errorf("errored") }

func (MockedWriter) Write(b []byte) (int, error) { return len(b), io.EOF }

func (MockedWriterFail) Write(b []byte) (int, error) { return len(b), fmt.Errorf("errored") }

func (MockedCloser) Close() error { return nil }

func (MockedSeeker) Seek(offset int64, whence int) (int64, error) { return 0, nil }

func (MockedSeekerFail) Seek(offset int64, whence int) (int64, error) {
	return 0, fmt.Errorf("errored")
}

// Check if a Close() call on FDs causes an error when at least one FD is nil.
func TestFDCloseFails(t *testing.T) {
	_, closer := NewFDs(nil, nil, nil, nil)
	if err := closer(); err == nil && !errors.Is(err, NilFDError) {
		t.Error("nil file descriptor should raise NilFDError")
	}
}

// Verify that a Close() call does not cause an error when all FDs are not nil.
func TestFDCloseOk(t *testing.T) {
	fds, _ := NewFDs(
		MockedReadCloser{},
		MockedWriteCloser{},
		MockedWriteCloser{},
		MockedReadWriteSeekCloser{},
	)
	err := fds.Close()
	if err != nil {
		t.Error("error should not be raised with open file descriptors")
	}
}

// Test if a Write() call does not bring about an error.
func TestDiscardWrite(t *testing.T) {
	buff := make([]byte, 512)
	if n, err := Discard.Write(buff); err != nil && n != len(buff) {
		t.Error("call to Write() on discard shoud not cause an error")
	}
}

// Test if a Close() call does not cause an error.
func TestDiscardClose(t *testing.T) {
	if err := Discard.Close(); err != nil {
		t.Error("call to Close() on Discard shoud never cause an error")
	}
}

// Verify if Wrap gets called successfully.
func TestWrapSuccess(t *testing.T) {
	cmd := MockedCmd{}
	fds, _ := NewFDs(
		MockedReadCloser{},
		MockedWriteCloser{},
		MockedWriteCloser{},
		MockedReadWriteSeekCloser{},
	)
	err := Wrap(cmd, fds)
	if err != nil {
		t.Error("Wrap() should pass without any errors in this scenario")
	}
}

// Test the constructor interface agreement.
func TestCmdConstructor(t *testing.T) {
	var _ Runner
	_ = Cmd("black", MockedWriter{}, MockedWriter{}, []string{"main.py"}...)
}

// Check if Wrap fails in the predefined scenarios.
func TestWrapRunFailed(t *testing.T) {
	cases := []struct {
		name string
		cmd  Runner
		fds  *FDs
	}{
		{
			"runner-fail",
			MockedCmdFail{},
			&FDs{
				MockedReadCloser{},
				MockedWriteCloser{},
				MockedWriteCloser{},
				MockedReadWriteSeekCloser{},
			},
		},
		{
			"stdout-fail",
			MockedCmd{},
			&FDs{
				MockedReadCloser{},
				MockedWriteCloserFail{},
				MockedWriteCloser{},
				MockedReadWriteSeekCloser{},
			},
		},
		{
			"tmpfile-write-fail",
			MockedCmd{},
			&FDs{
				MockedReadCloser{},
				MockedWriteCloser{},
				MockedWriteCloser{},
				MockedReadWriteSeekCloserFail2{},
			},
		},
		{
			"tmpfile-read-fail",
			MockedCmd{},
			&FDs{
				MockedReadCloser{},
				MockedWriteCloser{},
				MockedWriteCloser{},
				MockedReadWriteSeekCloserFail{},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := Wrap(c.cmd, c.fds)
			if err == nil {
				t.Error("the following Wrap attributes should make it fail")
			}
		})
	}
}
