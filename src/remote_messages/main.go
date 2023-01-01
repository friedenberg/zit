package remote_messages

import (
	"encoding/gob"
	"net"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type DialogueType int

const (
	DialogueTypeUnknown = DialogueType(iota)
	DialogueTypeDirector
	DialogueTypePull
	DialogueTypePullObjekten
	DialogueTypePullAkte
)

type Dialogue struct {
	typ   DialogueType
	conn  net.Conn
	stage *stage
	dec   *gob.Decoder
	enc   *gob.Encoder
}

func (s Dialogue) Write(p []byte) (n int, err error) {
	if n, err = s.conn.Write(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s Dialogue) Read(p []byte) (n int, err error) {
	if n, err = s.conn.Read(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s Dialogue) Close() (err error) {
	if err = s.conn.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s Dialogue) Type() DialogueType {
	return s.typ
}

func (s Dialogue) Send(e any) (err error) {
	if err = s.enc.Encode(e); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s Dialogue) Receive(e any) (err error) {
	if err = s.dec.Decode(e); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
