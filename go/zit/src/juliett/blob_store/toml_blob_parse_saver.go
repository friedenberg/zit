package blob_store

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/alfa/toml"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
)

type tomlBlobParseSaver[
	O interfaces.Blob[O],
	OPtr interfaces.BlobPtr[O],
] struct {
	awf              interfaces.BlobWriterFactory
	ignoreTomlErrors bool
}

func MakeTomlBlobParseSaver[
	O interfaces.Blob[O],
	OPtr interfaces.BlobPtr[O],
](awf interfaces.BlobWriterFactory,
) tomlBlobParseSaver[O, OPtr] {
	return tomlBlobParseSaver[O, OPtr]{
		awf: awf,
	}
}

func MakeTextParserIgnoreTomlErrors[
	O interfaces.Blob[O],
	OPtr interfaces.BlobPtr[O],
](awf interfaces.BlobWriterFactory,
) tomlBlobParseSaver[O, OPtr] {
	return tomlBlobParseSaver[O, OPtr]{
		awf:              awf,
		ignoreTomlErrors: true,
	}
}

func (f tomlBlobParseSaver[O, OPtr]) ParseBlob(
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

		ui.TodoP1("handle url parsing / validation")
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
