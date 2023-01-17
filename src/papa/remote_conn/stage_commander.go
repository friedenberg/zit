package remote_conn

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/india/erworben"
	"github.com/friedenberg/zit/src/oscar/umwelt"
)

type MessageHiCommander struct {
	DialogueType
	CliKonfig erworben.Cli
}

type StageCommander struct {
	remoteActorCmd      *exec.Cmd
	konfigCli           erworben.Cli
	wg                  *sync.WaitGroup
	chRemoteCommandDone chan struct{}
	stage
}

func (s StageCommander) ChanRemoteCommandDone() <-chan struct{} {
	return s.chRemoteCommandDone
}

func (s StageCommander) Close() (err error) {
	//TODO-P3 determine if this is the right place
	if err = s.MainDialogue().Close(); err != nil {
		if errors.IsErrno(err, syscall.EPIPE) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	s.wg.Wait()

	if err = s.remoteActorCmd.Wait(); err != nil {
		err = errors.Wrap(err)
		// errors.Err().Printf("close error: %s", err)
		return
	}

	close(s.chRemoteCommandDone)

	return
}

func (c StageCommander) ShouldIgnoreConnectionError(in error) (ok bool) {
	select {
	case <-c.chRemoteCommandDone:
		if errors.Is(in, net.ErrClosed) {
			ok = true
		}

	default:
	}

	return
}

func MakeStageCommander(
	u *umwelt.Umwelt,
	from string,
	command string,
) (s *StageCommander, err error) {
	s = &StageCommander{
		wg:                  &sync.WaitGroup{},
		konfigCli:           u.Konfig().Cli(),
		chRemoteCommandDone: make(chan struct{}),
	}

	s.remoteActorCmd = exec.Command(
		u.Standort().Executable(),
		"listen",
		"-dir-zit",
		from,
		command,
	)

	var rErr io.ReadCloser

	if rErr, err = s.remoteActorCmd.StderrPipe(); err != nil {
		err = errors.Wrap(err)
		return
	}

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

	copyWithPrefix := func(r *bufio.Reader, w io.Writer) {
		defer s.wg.Done()

		for {
			var line string

			if line, err = r.ReadString('\n'); err != nil {
				break
			}

			fmt.Fprintf(w, "remote: %s", line)
		}
	}

	s.wg.Add(2)
	go copyWithPrefix(rb, os.Stdout)
	go copyWithPrefix(bufio.NewReader(rErr), os.Stderr)

	s.sockPath = strings.TrimSpace(s.sockPath)

	if s.address, err = net.ResolveUnixAddr("unix", s.sockPath); err != nil {
		err = errors.Wrap(err)
		return
	}

	if s.mainDialogue, err = s.StartDialogue(
		DialogueTypeMain,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *StageCommander) StartDialogue(t DialogueType) (d Dialogue, err error) {
	if d, err = makeDialogueDial(&s.stage, t); err != nil {
		err = errors.Wrap(err)
		return
	}

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
