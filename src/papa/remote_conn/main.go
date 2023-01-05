package remote_conn

import (
	"encoding/gob"
	"net"

	"github.com/friedenberg/zit/src/alfa/errors"
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
	// errors.Log().Printf("%s sending %T:%v", s.Type(), e, e)
	// defer errors.Log().Printf("%s sent %T:%v", s.Type(), e, e)

	if err = s.enc.Encode(e); err != nil {
		if errors.IsEOF(err) {
			errors.Log().Caller(1, "%s EOF", s.Type())
		}

		err = errors.Wrapf(err, "%s", s.Type())
		return
	}

	return
}

func (s Dialogue) Receive(e any) (err error) {
	// errors.Log().Printf("%s receiving %T:%v", s.Type(), e, e)
	// defer errors.Log().Printf("%s received %T:%v", s.Type(), e, e)

	if err = s.dec.Decode(e); err != nil {
		if errors.IsEOF(err) {
			errors.Log().Caller(1, "%s EOF", s.Type())
		} else {
			err = errors.Wrapf(err, "%s", s.Type())
			return
		}
	}

	return
}
