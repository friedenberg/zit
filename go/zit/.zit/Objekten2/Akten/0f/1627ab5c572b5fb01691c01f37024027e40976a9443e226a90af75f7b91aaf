package remote_conn

import (
	"fmt"
	"io"
	"net"
	"path/filepath"
	"runtime/debug"
	"syscall"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type Listener interface {
	Listen() error
}

type MessageHiSoldier struct {
	Angeboren interfaces.ImmutableConfig
}

type SoldierDialogueChanElement struct {
	Dialogue
	MessageHiCommander
	error
}

type StageSoldier struct {
	Angeboren                 interfaces.ImmutableConfigGetter
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

func MakeStageSoldier(u *env.Local) (
	s *StageSoldier,
	err error,
) {
	s = &StageSoldier{
		Angeboren:                 u.GetConfig(),
		chStopWaitingForDialogues: make(chan struct{}),
		handlers:                  make(map[DialogueType]func(Dialogue) error),
	}

	var d string

	if d, err = u.GetDirectoryLayout().TempOS.DirTemp(); err != nil {
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
	u.GetConfig().SetCliFromCommander(el.MessageHiCommander.CliKonfig)
	ui.Log().Printf("set konfig")

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
	ui.Log().Printf("waiting for handlers")

	for {
		if !errMulti.Empty() {
			return
		}

		var el SoldierDialogueChanElement

		ui.Log().Printf("waiting for connection")

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
						errMulti.Add(errors.Wrapf(err, "%s", debug.Stack()))
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

			ui.Log().Printf("connection accepted")

			var h func(Dialogue) error
			ok := false

			if h, ok = s.handlers[el.Dialogue.Type()]; !ok {
				err := errors.Errorf(
					"unregistered dialogue type: %s", el.Dialogue.Type(),
				)

				errMulti.Add(err)
				return
			}

			ui.Log().Printf("found handler: %s", el.Dialogue.Type())

			ui.Log().Printf("start handler: %s", el.Dialogue.Type())
			defer ui.Log().Printf("end handler: %s", el.Dialogue.Type())

			if err := h(el.Dialogue); err != nil {
				errMulti.Add(errors.Wrap(err))
			}
		}()
	}
}

func (s *StageSoldier) AwaitDialogue() (out SoldierDialogueChanElement) {
	if out.Dialogue, out.MessageHiCommander, out.error = makeDialogueListen(
		s.Angeboren,
		&s.stage,
		s.listener,
	); out.error != nil {
		out.error = errors.Wrap(out.error)
		return
	}

	return
}
