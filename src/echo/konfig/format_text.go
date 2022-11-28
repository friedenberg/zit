package konfig

import (
	"bufio"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/toml"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/metadatei_io"
)

type FormatText struct {
	af metadatei_io.AkteIOFactory
}

func MakeFormatText(af metadatei_io.AkteIOFactory) *FormatText {
	return &FormatText{
		af: af,
	}
}

func (c *FormatText) ReadFormat(r1 io.Reader, k *Objekte) (n int64, err error) {
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

		//TODO fix issue with wrap not adding stack
		if err = td.Decode(&k.tomlKonfig); err != nil {
			if errors.IsEOF(err) {
				err = nil
			} else {
				pr.CloseWithError(errors.Wrap(err))
			}
		}

		if k.Akte, err = makeCompiled(k.tomlKonfig); err != nil {
			err = errors.Wrap(err)
			return
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

	k.Sha = aw.Sha()

	return
}
