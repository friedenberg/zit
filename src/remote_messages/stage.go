package remote_messages

import (
	"encoding/gob"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/oscar/umwelt"
)

type Stage struct {
	listener       net.Listener
	remoteActorCmd *exec.Cmd

	sockPath string
	conn     net.Conn
	dec      *gob.Decoder
	enc      *gob.Encoder
}

func (s Stage) Close() (err error) {
	if s.remoteActorCmd != nil {
		if err = s.remoteActorCmd.Wait(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if s.listener != nil {
		if err = syscall.Unlink(s.sockPath); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func MakeStageSender(u *umwelt.Umwelt, from string) (s *Stage, err error) {
	s = &Stage{}

	var d string

	if d, err = os.MkdirTemp("", ""); err != nil {
		err = errors.Wrap(err)
		return
	}

	s.sockPath = filepath.Join(d, "zit-pull.sock")

	s.remoteActorCmd = exec.Command(
		u.Standort().Executable(),
		"listen",
		"-dir-zit",
		from,
		s.sockPath,
	)

	s.remoteActorCmd.Stderr = os.Stderr

	var r io.ReadCloser

	if r, err = s.remoteActorCmd.StdoutPipe(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.remoteActorCmd.Start(); err != nil {
		err = errors.Wrap(err)
		return
	}

	expected := "listening\n"
	actual := &strings.Builder{}

	if _, err = io.CopyN(actual, r, int64(len([]byte(expected)))); err != nil {
		err = errors.Wrap(err)
		return
	}

	if expected != actual.String() {
		err = errors.Errorf(
			"expected listener to emit %q but got %q",
			expected,
			actual.String(),
		)

		return
	}

	if s.conn, err = net.Dial("unix", s.sockPath); err != nil {
		err = errors.Wrap(err)
		return
	}

	s.enc = gob.NewEncoder(s.conn)
	s.dec = gob.NewDecoder(s.conn)

	return
}

func MakeStageReceiver(u *umwelt.Umwelt, sockPath string) (s *Stage, err error) {
	s = &Stage{
		sockPath: sockPath,
	}

	if s.listener, err = net.Listen("unix", sockPath); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, err = io.WriteString(u.Out(), "listening\n"); err != nil {
		err = errors.Wrap(err)
		return
	}

	if s.conn, err = s.listener.Accept(); err != nil {
		err = errors.Wrap(err)
		return
	}

	s.enc = gob.NewEncoder(s.conn)
	s.dec = gob.NewDecoder(s.conn)

	return
}

func (s *Stage) Send(e any) (err error) {
	if err = s.enc.Encode(e); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Stage) Receive(e any) (err error) {
	if err = s.dec.Decode(e); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
