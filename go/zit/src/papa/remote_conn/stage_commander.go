package remote_conn

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/golf/mutable_config_blobs"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type MessageHiCommander struct {
	DialogueType
	CliKonfig mutable_config_blobs.Cli
	Angeboren interfaces.ImmutableConfig
}

type StageCommander struct {
	Angeboren           interfaces.ImmutableConfigGetter
	remoteActorCmd      *exec.Cmd
	konfigCli           mutable_config_blobs.Cli
	wg                  *sync.WaitGroup
	chRemoteCommandDone chan struct{}
	stage
}

func (s StageCommander) ChanRemoteCommandDone() <-chan struct{} {
	return s.chRemoteCommandDone
}

func (s StageCommander) Close() (err error) {
	// TODO-P3 determine if this is the right place
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
	u *env.Local,
	from string,
	command string,
) (s *StageCommander, err error) {
	if from == "" {
		err = errors.Errorf("empty from")
		return
	}

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs

		if s.MainDialogue().conn != nil {
			s.MainDialogue().Close()
		}
	}()

	s = &StageCommander{
		Angeboren:           u.GetConfig(),
		wg:                  &sync.WaitGroup{},
		konfigCli:           u.GetConfig().Cli(),
		chRemoteCommandDone: make(chan struct{}),
	}

	s.remoteActorCmd = exec.Command(
		u.GetDirectoryLayout().Executable(),
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

	// TODO-P2 make it possible to output and check for path simulataneously
	if s.sockPath, err = rb.ReadString('\n'); err != nil {
		err = errors.Wrapf(err, "Cmd: %s", s.remoteActorCmd.String())
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
		Angeboren:    s.Angeboren.GetImmutableConfig(),
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

	d.Angeboren = msgTheirHi.Angeboren

	return
}
