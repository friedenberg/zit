package typ

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/toml"
	"github.com/friedenberg/zit/src/bravo/metadatei_io"
	"github.com/friedenberg/zit/src/charlie/sha"
)

type FormatText struct {
	arf             metadatei_io.AkteIOFactory
	IgnoreTypErrors bool
}

func MakeFormatText(arf metadatei_io.AkteIOFactory) *FormatText {
	return &FormatText{
		arf: arf,
	}
}

func (f FormatText) ReadFormat(r io.Reader, t *Typ) (n int64, err error) {
	var aw sha.WriteCloser

	if aw, err = f.arf.AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, aw.Close)

	pr, pw := io.Pipe()

	chDone := make(chan struct{})

	go func() {
		td := toml.NewDecoder(pr)

		if err := td.Decode(t); err != nil {
			pr.CloseWithError(err)
		}

		chDone <- struct{}{}
	}()

	mw := io.MultiWriter(
		aw,
		pw,
	)

	if n, err = io.Copy(mw, r); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f FormatText) WriteFormat(w io.Writer, t *Typ) (n int64, err error) {
	var ar sha.ReadCloser

	if ar, err = f.arf.AkteReader(t.Akte.Sha); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, ar.Close)

	mw := metadatei_io.Writer{
		// Metadatei: ,
		Akte: ar,
	}

	if n, err = mw.WriteTo(w); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
