package remote_conn

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"syscall"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type Listener interface {
	Listen() error
}

type MessageHiSoldier struct{}

type SoldierDialogueChanElement struct {
	Dialogue
	MessageHiCommander
	error
}

type StageSoldier struct {
	listener                  *net.UnixListener
	chStopWaitingForDialogues chan struct{}
	chDialogue                chan SoldierDialogueChanElement
	handlers                  map[DialogueType]func(Dialogue) error
	stage
}

func (s StageSoldier) Close() (err error) {
	if s.listener != nil {
		if err = syscall.Unlink(s.sockPath); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func MakeStageSoldier(u *umwelt.Umwelt) (
	s *StageSoldier,
	err error,
) {
	s = &StageSoldier{
		chStopWaitingForDialogues: make(chan struct{}),
		handlers:                  make(map[DialogueType]func(Dialogue) error),
	}

	var d string

	if d, err = os.MkdirTemp("", ""); err != nil {
		err = errors.Wrap(err)
		return
	}

	s.sockPath = filepath.Join(d, "zit.sock")

	if s.address, err = net.ResolveUnixAddr("unix", s.sockPath); err != nil {
		err = errors.Wrap(err)
		return
	}

	if s.listener, err = net.ListenUnix("unix", s.address); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, err = io.WriteString(
		u.Out(),
		fmt.Sprintf("%s\n", s.sockPath),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var el SoldierDialogueChanElement

	if el = s.AwaitDialogue(); err != nil {
		err = errors.Wrap(err)
		return
	}

	s.mainDialogue = el.Dialogue
	u.KonfigPtr().SetCliFromCommander(el.MessageHiCommander.CliKonfig)
	errors.Log().Printf("set konfig")

	if err = u.Reset(); err != nil {
		err = errors.Wrap(err)
		return
	}

	err = el.error

	return
}

func (s *StageSoldier) RegisterHandler(
	t DialogueType,
	h func(Dialogue) error,
) {
	s.handlers[t] = h
}

func (s *StageSoldier) Listen() (err error) {
	errMulti := errors.MakeMulti()

	defer func() {
		errMulti.Add(err)

		if !errMulti.Empty() {
			err = errMulti
		}
	}()

	defer errors.Deferred(&err, s.Close)

	go func() {
		<-errMulti.ChanOnErr()
		errMulti.Add(s.MainDialogue().Close())
	}()

	go s.awaitRegisteredDialogueHandlers(errMulti)

	var done interface{}

	if err = s.MainDialogue().Receive(done); err != nil {
		if errors.IsEOF(err) {
			err = nil
		} else {
			errMulti.Add(errors.Wrap(err))
			return
		}
	}

	return
}

func (s *StageSoldier) awaitRegisteredDialogueHandlers(errMulti errors.Multi) {
	errors.Log().Printf("waiting for handlers")

	for {
		if !errMulti.Empty() {
			return
		}

		var el SoldierDialogueChanElement

		errors.Log().Printf("waiting for connection")

		if el = s.AwaitDialogue(); el.error != nil {
			errMulti.Add(el.error)
			return
		}

		if !errMulti.Empty() {
			return
		}

		go func() {
			defer func() {
				if e := recover(); e != nil {
					if err, ok := e.(error); ok {
						errMulti.Add(err)
					} else {
						panic(e)
					}
				}
			}()

			if !errMulti.Empty() {
				return
			}

			if el.Dialogue.Type() == DialogueTypeMain {
				err := errors.Errorf("receive request for main dialog after handshake")
				errMulti.Add(err)
				return
			}

			errors.Log().Printf("connection accepted")

			var h func(Dialogue) error
			ok := false

			if h, ok = s.handlers[el.Dialogue.Type()]; !ok {
				err := errors.Errorf(
					"unregistered dialogue type: %s", el.Dialogue.Type(),
				)

				errMulti.Add(err)
				return
			}

			errors.Log().Printf("found handler: %s", el.Dialogue.Type())

			errors.Log().Printf("start handler: %s", el.Dialogue.Type())
			defer errors.Log().Printf("end handler: %s", el.Dialogue.Type())

			if err := h(el.Dialogue); err != nil {
				errMulti.Add(errors.Wrap(err))
			}
		}()
	}
}

func (s *StageSoldier) AwaitDialogue() (out SoldierDialogueChanElement) {
	if out.Dialogue, out.MessageHiCommander, out.error = makeDialogueListen(
		&s.stage,
		s.listener,
	); out.error != nil {
		out.error = errors.Wrap(out.error)
		return
	}

	return
}
