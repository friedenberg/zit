package erworben

import (
	"bufio"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/alfa/toml"
	"github.com/friedenberg/zit/src/bravo/sha"
)

// TODO-P1 rename to TextFormat
type FormatText struct {
	af schnittstellen.AkteIOFactory
}

func MakeFormatText(af schnittstellen.AkteIOFactory) *FormatText {
	return &FormatText{
		af: af,
	}
}

func (c *FormatText) Parse(r1 io.Reader, k *Objekte) (n int64, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("toml unmarshalling panicked: %q", r)
		}
	}()

	r := bufio.NewReader(r1)

	pr, pw := io.Pipe()

	chDone := make(chan struct{})

	go func() {
		defer func() {
			close(chDone)
		}()

		td := toml.NewDecoder(pr)
		td.DisallowUnknownFields()

		//TODO-P4 fix issue with wrap not adding stack
		if err = td.Decode(&k.Akte); err != nil {
			if errors.IsEOF(err) {
				err = nil
			} else {
				pr.CloseWithError(errors.Wrap(err))
			}
		}
	}()

	var aw sha.WriteCloser

	if aw, err = c.af.AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, aw.Close)

	mw := io.MultiWriter(aw, pw)

	if n, err = io.Copy(mw, r); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = pw.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	<-chDone

	k.Sha = sha.Make(aw.Sha())

	return
}

func (f FormatText) Format(w io.Writer, t *Objekte) (n int64, err error) {
	var ar sha.ReadCloser

	if ar, err = f.af.AkteReader(t.Sha); err != nil {
		if errors.IsNotExist(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	defer errors.Deferred(&err, ar.Close)

	if n, err = io.Copy(w, ar); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
