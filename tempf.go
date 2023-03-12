package duct

import (
	"os"
)

const pattern = `duct-*`

type Tempf struct {
	*os.File
}

func NewTempf() (*Tempf, error) {
	f, err := os.CreateTemp("", pattern)
	if err != nil {
		return nil, err
	}
	t := Tempf{f}
	return &t, nil
}
