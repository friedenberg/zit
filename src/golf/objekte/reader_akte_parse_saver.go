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
	parser AkteParseSaver[OPtr]
}

func MakeReaderAkteParseSaver[
	O Objekte[O],
	OPtr ObjektePtr[O],
](
	awf schnittstellen.AkteWriterFactory,
	parser AkteParseSaver[OPtr],
) readerAkteParseSaver[O, OPtr] {
	return readerAkteParseSaver[O, OPtr]{
		awf:    awf,
		parser: parser,
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

	var (
		n1  int64
		sh1 schnittstellen.Sha
	)

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

		if sh1, n1, err = f.parser.ParseSaveAkte(pr, t); err != nil {
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

	if n != n1 {
		err = errors.Errorf(
			"parser read %d bytes while akte writer read %d",
			n1,
			n,
		)

		return
	}

	if err = <-chDone; err != nil {
		err = errors.Wrap(err)
		return
	}

	sh = aw.Sha()

	if !sh.EqualsSha(sh1) {
		err = errors.Errorf(
			"parser read %s while akte writer read %s",
			sh1,
			sh,
		)

		return
	}

	return
}
