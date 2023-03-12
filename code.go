package duct

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

/*
1. Take the standard input from the program.
  - The reader might take io.Reader interface

2. Save it to a temporary file.
  - Use tempfile

3. Pass the file name to the shell command call.
4. Pass stdout and stderr to /dev/null.
5. Read text from the temporary files when the command finishes.
6. Capture stderr to report an error
*/
func Do() {
	in := bufio.NewReader(os.Stdin)
	f, err := NewTempf()
	in.WriteTo(f)
	defer f.Close()
	defer os.Remove(f.Name())
	if err != nil {
		return
	}
	f.Seek(0, 0)
	if err != nil {
		return
	}
	f.Name()
	cmdName := "black"
	cmd := exec.Command(cmdName, f.Name())
	cmd.Stdout = ioutil.Discard
	cmd.Stderr = ioutil.Discard
	err = cmd.Run()
	buff := make([]byte, 512)
	out := []byte{}
	for {
		_, err := f.Read(buff)
		if err == io.EOF {
			break
		}
		out = append(out, buff...)

	}
	fmt.Println(string(out))
}
