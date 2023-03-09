package duct

import (
	"io"
	"os"
	"sync"
)

type ReadWriteSeekCloser interface {
	io.ReadSeekCloser
	io.Writer
}

type Tempf struct {
	f *os.File
	sync.RWMutex
	ReadWriteSeekCloser
}

func NewTempf() *Tempf {
	t := Tempf{f: &os.File{}}
	return &t
}
