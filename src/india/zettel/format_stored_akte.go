package zettel

import (
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/toml"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/sha"
)

type FormatStoredAkte struct {
	arf  gattung.AkteIOFactory
	toml bool
}

func MakeFormatStoredAkte(arf gattung.AkteIOFactory) *FormatStoredAkte {
	return &FormatStoredAkte{
		arf: arf,
	}
}

func (f FormatStoredAkte) ReadFormat(r io.Reader, t *Stored) (n int64, err error) {
	var aw sha.WriteCloser

	if aw, err = f.arf.AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, aw.Close)

	var pw io.WriteCloser
	var pr *io.PipeReader

	chDone := make(chan struct{})

	if f.toml {
		pr, pw = io.Pipe()

		go func() {
			defer func() {
				close(chDone)
			}()

			td := toml.NewDecoder(pr)

			var a map[string]interface{}

			if err := td.Decode(&a); err != nil {
				if !errors.IsEOF(err) {
					pr.CloseWithError(err)
				}
			}
		}()
	} else {
		pw = files.NopWriteCloser{Writer: io.Discard}
		close(chDone)
	}

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

	t.Objekte.Akte = aw.Sha()

	return
}

func (f FormatStoredAkte) WriteFormat(w io.Writer, t *Stored) (n int64, err error) {
	var ar sha.ReadCloser

	if ar, err = f.arf.AkteReader(t.Objekte.Akte); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, ar.Close)

	if _, err = io.WriteString(w, fmt.Sprintf("[%s]", t)); err != nil {
		err = errors.Wrap(err)
		return
	}

	if n, err = io.Copy(w, ar); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
