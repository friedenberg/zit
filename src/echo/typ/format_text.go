package typ

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/toml"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/metadatei_io"
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

func (f FormatText) ReadFormat(r io.Reader, t *Akte) (n int64, err error) {
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

		if err := td.Decode(&t.KonfigTyp); err != nil {
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

	t.Sha = aw.Sha()

	return
}

func (f FormatText) WriteFormat(w io.Writer, t *Akte) (n int64, err error) {
	var ar sha.ReadCloser

	if ar, err = f.arf.AkteReader(t.Sha); err != nil {
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
