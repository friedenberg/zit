package remote_messages

import (
	"bufio"
	"encoding/gob"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/india/konfig"
	"github.com/friedenberg/zit/src/oscar/umwelt"
)

type MessageHiCommander struct {
	DialogueType
	CliKonfig konfig.Cli
}

type StageCommander struct {
	remoteActorCmd *exec.Cmd
	konfigCli      konfig.Cli
	stage
}

func (s StageCommander) Close() (err error) {
	if err = s.remoteActorCmd.Wait(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func MakeStageCommander(
	u *umwelt.Umwelt,
	from string,
) (s *StageCommander, err error) {
	s = &StageCommander{
		konfigCli: u.Konfig().Cli(),
	}

	s.remoteActorCmd = exec.Command(
		u.Standort().Executable(),
		"listen",
		"-dir-zit",
		from,
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

	rb := bufio.NewReader(r)

	if s.sockPath, err = rb.ReadString('\n'); err != nil {
		err = errors.Wrap(err)
		return
	}

	go func() {
		io.Copy(os.Stdout, rb)
	}()

	s.sockPath = strings.TrimSpace(s.sockPath)

	if s.mainDialogue, err = s.StartDialogue(
		DialogueTypeDirector,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *StageCommander) StartDialogue(t DialogueType) (d Dialogue, err error) {
	d.typ = t
	d.stage = &s.stage

	if d.conn, err = net.Dial("unix", s.sockPath); err != nil {
		err = errors.Wrap(err)
		return
	}

	d.enc = gob.NewEncoder(d.conn)
	d.dec = gob.NewDecoder(d.conn)

	msgOurHi := MessageHiCommander{
		DialogueType: d.Type(),
		CliKonfig:    s.konfigCli,
	}

	if err = d.Send(msgOurHi); err != nil {
		err = errors.Wrap(err)
		return
	}

	var msgTheirHi MessageHiSoldier

	if err = d.Receive(&msgTheirHi); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
