package script_value

import (
	"io"
	"os"
	"os/exec"

	"github.com/friedenberg/zit/alfa/logz"
	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/bravo/open_file_guard"
)

type ScriptValue struct {
	script string
	cmd    *exec.Cmd
	file   *os.File
}

func (s ScriptValue) String() string {
	return s.script
}

func (s ScriptValue) IsEmpty() bool {
	return s.script == ""
}

func (s *ScriptValue) Set(v string) (err error) {
	s.script = v

	return
}

func (s *ScriptValue) Run(input string) (r io.Reader, err error) {
	if s.IsEmpty() {
		if input == "" || input == "-" {
			r = os.Stdin
		} else {
			if s.file, err = open_file_guard.Open(input); err != nil {
				err = errors.Error(err)
				return
			}

			r = s.file
		}

		return
	}

	if input == "" || input == "-" {
		s.cmd = exec.Command(s.script)
		s.cmd.Stdin = os.Stdin
	} else {
		s.cmd = exec.Command(s.script, input)
	}

	if r, err = s.cmd.StdoutPipe(); err != nil {
		logz.Fatal(err)
		return
	}

	logz.Print("starting")
	s.cmd.Start()

	return
}

func (s *ScriptValue) Close() (err error) {
	logz.Print()

	if s.file != nil {
		err = open_file_guard.Close(s.file)
	}

	if s.cmd != nil {
		err = s.cmd.Wait()
	}

	return
}
