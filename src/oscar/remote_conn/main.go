package remote_conn

import (
	"encoding/gob"
	"net"
	"syscall"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type Dialogue struct {
	typ   DialogueType
	conn  *net.UnixConn
	stage *stage
	dec   *gob.Decoder
	enc   *gob.Encoder
}

func makeDialogueListen(
	s *stage,
	l *net.UnixListener,
) (d Dialogue, msg MessageHiCommander, err error) {
	d.stage = s

	if d.conn, err = l.AcceptUnix(); err != nil {
		//TODO-P2 determine what errors accept can throw
		if errors.IsErrno(err, syscall.ENODATA) || true {
			panic(errors.Wrapf(err, "ErrorExact: %#v", err))
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	d.enc = gob.NewEncoder(d.conn)
	d.dec = gob.NewDecoder(d.conn)

	msgOurHi := MessageHiSoldier{}

	if err = d.Send(msgOurHi); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = d.Receive(
		&msg,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	d.typ = msg.DialogueType

	return
}

func makeDialogueDial(
	s *stage,
	t DialogueType,
) (d Dialogue, err error) {
	d.stage = s
	d.typ = t

	for {
		//TODO-P2 timeout
		if d.conn, err = net.DialUnix("unix", nil, s.address); err != nil {
			//TODO-P5 why is the hex 0x3d which is ENODATA in docs?
			if errors.IsErrno(err, syscall.ECONNREFUSED) {
				WaitForConnectionLicense()
				continue
			} else {
				err = errors.Wrap(err)
				return
			}
		}

		break
	}

	d.enc = gob.NewEncoder(d.conn)
	d.dec = gob.NewDecoder(d.conn)

	return
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
	ReturnConnLicense()

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
