package duct

import (
	"bufio"
)

func Wrap(name string, fds *Fds) error {
	in := bufio.NewReader(fds.Stdin)
	_, err := in.WriteTo(fds.Tempf)
	if err != nil {
		return err
	}
	_, err = fds.Tempf.Seek(0, 0)
	if err != nil {
		return err
	}
	cmd := NewFormatter(name, fds.Stdout, fds.Stderr, fds.Tempf.Name())
	err = cmd.Run()
	if err != nil {
		return err
	}
	out := bufio.NewWriter(fds.Stdout)
	_, err = out.ReadFrom(fds.Tempf)
	if err != nil {
		return err
	}
	return nil
}
