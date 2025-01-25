package typed_blob_store

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/alfa/toml"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
)

type tomlBlobDecoder[
	O any,
	OPtr interfaces.Ptr[O],
] struct {
	awf              interfaces.BlobWriter
	ignoreTomlErrors bool
}

func MakeTomlBlobDecoderSaver[
	O any,
	OPtr interfaces.Ptr[O],
](
	blobWriter interfaces.BlobWriter,
) tomlBlobDecoder[O, OPtr] {
	return tomlBlobDecoder[O, OPtr]{
		awf: blobWriter,
	}
}

func MakeTomlDecoderIgnoreTomlErrors[
	O any,
	OPtr interfaces.Ptr[O],
](awf interfaces.BlobWriter,
) tomlBlobDecoder[O, OPtr] {
	return tomlBlobDecoder[O, OPtr]{
		awf:              awf,
		ignoreTomlErrors: true,
	}
}

func (f tomlBlobDecoder[O, OPtr]) DecodeFrom(
	t OPtr,
	r io.Reader,
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
