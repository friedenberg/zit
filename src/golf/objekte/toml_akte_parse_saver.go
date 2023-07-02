package objekte

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/alfa/toml"
	"github.com/friedenberg/zit/src/bravo/sha"
)

type tomlAkteParseSaver[
	O Akte[O],
	OPtr AktePtr[O],
] struct {
	awf              schnittstellen.AkteWriterFactory
	ignoreTomlErrors bool
}

func MakeTomlAkteParseSaver[
	O Akte[O],
	OPtr AktePtr[O],
](awf schnittstellen.AkteWriterFactory,
) tomlAkteParseSaver[O, OPtr] {
	return tomlAkteParseSaver[O, OPtr]{
		awf: awf,
	}
}

func MakeTextParserIgnoreTomlErrors[
	O Akte[O],
	OPtr AktePtr[O],
](awf schnittstellen.AkteWriterFactory,
) tomlAkteParseSaver[O, OPtr] {
	return tomlAkteParseSaver[O, OPtr]{
		awf:              awf,
		ignoreTomlErrors: true,
	}
}

func (f tomlAkteParseSaver[O, OPtr]) ParseSaveAkte(
	r io.Reader,
	t OPtr,
) (sh schnittstellen.ShaLike, n int64, err error) {
	var aw sha.WriteCloser

	if aw, err = f.awf.AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, aw)

	pr, pw := io.Pipe()
	td := toml.NewDecoder(pr)

	chDone := make(chan error)

	go func(pr *io.PipeReader) {
		var err error
		defer func() {
			chDone <- err
			close(chDone)
		}()

		defer func() {
			if r := recover(); r != nil {
				if f.ignoreTomlErrors {
					err = nil
				} else {
					err = toml.MakeError(errors.Errorf("panicked during toml decoding: %s", r))
					pr.CloseWithError(errors.Wrap(err))
				}
			}
		}()

		if err = td.Decode(t); err != nil {
			switch {
			case !errors.IsEOF(err) && !f.ignoreTomlErrors:
				err = errors.Wrap(toml.MakeError(err))
				pr.CloseWithError(err)

			case !errors.IsEOF(err) && f.ignoreTomlErrors:
				err = nil
			}
		}

		errors.TodoP1("handle url parsing / validation")
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

	sh = aw.GetShaLike()

	return
}
