package akten

import (
	"io"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/alfa/toml"
)

type tomlAkteParseSaver[
	O schnittstellen.Akte[O],
	OPtr schnittstellen.AktePtr[O],
] struct {
	awf              schnittstellen.AkteWriterFactory
	ignoreTomlErrors bool
}

func MakeTomlAkteParseSaver[
	O schnittstellen.Akte[O],
	OPtr schnittstellen.AktePtr[O],
](awf schnittstellen.AkteWriterFactory,
) tomlAkteParseSaver[O, OPtr] {
	return tomlAkteParseSaver[O, OPtr]{
		awf: awf,
	}
}

func MakeTextParserIgnoreTomlErrors[
	O schnittstellen.Akte[O],
	OPtr schnittstellen.AktePtr[O],
](awf schnittstellen.AkteWriterFactory,
) tomlAkteParseSaver[O, OPtr] {
	return tomlAkteParseSaver[O, OPtr]{
		awf:              awf,
		ignoreTomlErrors: true,
	}
}

func (f tomlAkteParseSaver[O, OPtr]) ParseAkte(
	r io.Reader,
	t OPtr,
) (n int64, err error) {
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

	if n, err = io.Copy(pw, r); err != nil {
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

	return
}
