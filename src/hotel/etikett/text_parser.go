package etikett

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/alfa/toml"
	"github.com/friedenberg/zit/src/bravo/sha"
)

// TODO-P1 rename to TextFormat
type textParser struct {
	arf schnittstellen.AkteIOFactory
}

func MakeTextParser(arf schnittstellen.AkteIOFactory) textParser {
	return textParser{
		arf: arf,
	}
}

func (f textParser) ParseSaveAkte(
	r io.Reader,
	t *Objekte,
) (sh schnittstellen.Sha, n int64, err error) {
	var aw sha.WriteCloser

	if aw, err = f.arf.AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, aw)

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

	sh = sha.Make(aw.Sha())

	return
}

func (f textParser) Parse(r io.Reader, t *Objekte) (n int64, err error) {
	var sh schnittstellen.Sha

	sh, n, err = f.ParseSaveAkte(r, t)

	t.Sha = sha.Make(sh)

	return
}
