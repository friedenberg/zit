package etikett

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/alfa/toml"
	"github.com/friedenberg/zit/src/bravo/sha"
)

// TODO-P1 rename to TextFormat
type FormatText struct {
	arf schnittstellen.AkteIOFactory
}

func MakeFormatText(arf schnittstellen.AkteIOFactory) *FormatText {
	return &FormatText{
		arf: arf,
	}
}

func (f FormatText) Parse(r io.Reader, t *Objekte) (n int64, err error) {
	var aw sha.WriteCloser

	if aw, err = f.arf.AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, aw.Close)

	pr, pw := io.Pipe()
	td := toml.NewDecoder(pr)

	chDone := make(chan struct{})

	go func() {
		defer func() {
			close(chDone)
		}()

		if err := td.Decode(&t.Akte); err != nil {
			if !errors.IsEOF(err) {
				pr.CloseWithError(err)
			}
		}
	}()

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

	t.Sha = sha.Make(aw.Sha())

	return
}

func (f FormatText) Format(w io.Writer, t *Objekte) (n int64, err error) {
	var ar sha.ReadCloser

	if ar, err = f.arf.AkteReader(t.Sha); err != nil {
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