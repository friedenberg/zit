package remote_messages

import (
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/oscar/umwelt"
)

type MessageHiSoldier struct{}

type SoldierDialogueChanElement struct {
	Dialogue
	MessageHiCommander
	error
}

type StageSoldier struct {
	listener                  net.Listener
	wg                        *sync.WaitGroup
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

	select {
	case <-s.chStopWaitingForDialogues:
	default:
		close(s.chStopWaitingForDialogues)
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
		wg:                        &sync.WaitGroup{},
	}

	var d string

	if d, err = os.MkdirTemp("", ""); err != nil {
		err = errors.Wrap(err)
		return
	}

	s.sockPath = filepath.Join(d, "zit.sock")

	if s.listener, err = net.Listen("unix", s.sockPath); err != nil {
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

func (s *StageSoldier) RegisterHandlerWithUmwelt(
	t DialogueType,
	u *umwelt.Umwelt,
	h func(*umwelt.Umwelt, Dialogue) error,
) {
	s.handlers[t] = func(d Dialogue) (err error) {
		return h(u, d)
	}
}

func (s *StageSoldier) RegisterHandler(
	t DialogueType,
	h func(Dialogue) error,
) {
	s.handlers[t] = h
}

func (s *StageSoldier) Listen() (err error) {
	defer errors.Deferred(&err, s.Close)

	chErr := make(chan error)

	go func() {
		for err1 := range chErr {
			err = errors.MakeMulti(err, err1)

			select {
			case <-s.chStopWaitingForDialogues:
			default:
				close(s.chStopWaitingForDialogues)
			}
		}
	}()

	go s.awaitRegisteredDialogueHandlers(chErr)

	s.wg.Wait()

	if err = s.MainDialogue().Receive(&MessageDone{}); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *StageSoldier) awaitRegisteredDialogueHandlers(chErr chan<- error) {
	errors.Log().Printf("waiting for handlers")

	for {
		select {
		case <-s.chStopWaitingForDialogues:
			//TODO-P2 handle remaining connections
			return

		default:
			var el SoldierDialogueChanElement

			errors.Log().Printf("waiting for connection")

			if el = s.AwaitDialogue(); el.error != nil {
				chErr <- errors.Wrap(el.error)
				continue
			}

			errors.Log().Printf("connection accepted")

			var h func(Dialogue) error
			ok := false

			if h, ok = s.handlers[el.Dialogue.Type()]; !ok {
				chErr <- errors.Errorf(
					"unregistered dialogue type: %d", el.Dialogue.Type(),
				)

				continue
			}

			errors.Log().Printf("found handler: %d", el.Dialogue.Type())

			s.wg.Add(1)

			go func() {
				defer s.wg.Done()
				errors.Log().Printf("start handler: %d", el.Dialogue.Type())
				defer errors.Log().Printf("end handler: %d", el.Dialogue.Type())

				if err := h(el.Dialogue); err != nil {
					chErr <- errors.Wrap(err)
				}
			}()
		}
	}
}

func (s *StageSoldier) AwaitDialogue() (out SoldierDialogueChanElement) {
	out.Dialogue.stage = &s.stage

	if out.Dialogue.conn, out.error = s.listener.Accept(); out.error != nil {
		out.error = errors.Wrap(out.error)
		return
	}

	out.Dialogue.enc = gob.NewEncoder(out.Dialogue.conn)
	out.Dialogue.dec = gob.NewDecoder(out.Dialogue.conn)

	msgOurHi := MessageHiSoldier{}

	if out.error = out.Dialogue.Send(msgOurHi); out.error != nil {
		out.error = errors.Wrap(out.error)
		return
	}

	if out.error = out.Dialogue.Receive(
		&out.MessageHiCommander,
	); out.error != nil {
		out.error = errors.Wrap(out.error)
		return
	}

	out.Dialogue.typ = out.MessageHiCommander.DialogueType

	return
}
