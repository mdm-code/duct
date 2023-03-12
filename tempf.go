package duct

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

const pattern = `duct-*`

type Tempf struct {
	*os.File
}

type ReadWriteSeekNamer interface {
	io.ReadWriteSeeker
	Name() string
}

func NewTempf() (*Tempf, error) {
	f, err := os.CreateTemp("", pattern)
	if err != nil {
		return nil, err
	}
	t := Tempf{f}
	return &t, nil
}

type Fds struct {
	Stdin          io.Reader
	Stdout, Stderr io.Writer
	Tempf          ReadWriteSeekNamer
}

// new fds with a closer function

func NewFormatter(name string, stdout, stderr io.Writer, files ...string) *exec.Cmd {
	var err error
	var path string
	if filepath.Base(name) == name {
		path, err = exec.LookPath(name)
		if path != "" {
			name = path
		}
	}
	cmd := exec.Cmd{
		Path:   name,
		Args:   append([]string{name}, files...),
		Stdout: stdout,
		Stderr: stderr,
		Err:    err,
	}
	return &cmd
}
