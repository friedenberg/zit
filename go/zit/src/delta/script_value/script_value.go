package script_value

import (
	"io"
	"os"
	"os/exec"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/charlie/files"
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

func (s ScriptValue) Cmd() *exec.Cmd {
	return s.cmd
}

func (s *ScriptValue) RunWithInput() (w io.WriteCloser, r io.Reader, err error) {
	if s.IsEmpty() {
		err = errors.Errorf("empty script")
		return
	}

	s.cmd = exec.Command(s.script)

	if w, err = s.cmd.StdinPipe(); err != nil {
		errors.Wrap(err)
		return
	}

	if r, err = s.cmd.StdoutPipe(); err != nil {
		errors.Wrap(err)
		return
	}

	return
}

func (s *ScriptValue) Run(input string) (r io.Reader, err error) {
	if s.IsEmpty() {
		if input == "" || input == "-" {
			r = os.Stdin
		} else {
			if s.file, err = files.Open(input); err != nil {
				err = errors.Wrap(err)
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
		errors.Wrap(err)
		return
	}

	errors.Log().Print("starting")
	s.cmd.Start()

	return
}

func (s *ScriptValue) Close() (err error) {
	errors.Log().Print("closing script")
	defer errors.Log().Print("done closing script")

	if s.file != nil {
		errors.Log().Print("closing file")
		err = files.Close(s.file)
	}

	if s.cmd != nil {
		errors.Log().Print("waiting for script")
		err = s.cmd.Wait()
	}

	return
}
