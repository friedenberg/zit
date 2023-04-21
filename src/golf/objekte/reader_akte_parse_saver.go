package objekte

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/sha"
)

type readerAkteParseSaver[
	O Objekte[O],
	OPtr ObjektePtr[O],
] struct {
	awf    schnittstellen.AkteWriterFactory
	reader schnittstellen.Parser[O, OPtr]
}

func MakeReaderAkteParseSaver[
	O Objekte[O],
	OPtr ObjektePtr[O],
](
	awf schnittstellen.AkteWriterFactory,
	reader schnittstellen.Parser[O, OPtr],
) readerAkteParseSaver[O, OPtr] {
	return readerAkteParseSaver[O, OPtr]{
		awf:    awf,
		reader: reader,
	}
}

func (f readerAkteParseSaver[O, OPtr]) ParseSaveAkte(
	r io.Reader,
	t OPtr,
) (sh schnittstellen.Sha, n int64, err error) {
	var aw sha.WriteCloser

	if aw, err = f.awf.AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, aw)

	pr, pw := io.Pipe()

	chDone := make(chan error)

	go func(pr *io.PipeReader) {
		var err error
		defer func() {
			chDone <- err
			close(chDone)
		}()

		defer func() {
			if r := recover(); r != nil {
				err = errors.Errorf("panicked: %s", r)
				pr.CloseWithError(err)
			}
		}()

		if _, err = f.reader.Parse(pr, t); err != nil {
			pr.CloseWithError(err)
		}
	}(pr)

	mw := io.MultiWriter(aw, pw)

	if n, err = io.Copy(mw, r); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = pw.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = <-chDone; err != nil {
		err = errors.Wrap(err)
		return
	}

	sh = aw.Sha()

	return
}
