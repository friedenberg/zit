package akten

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
)

type readerAkteParseSaver[
	O interfaces.Blob[O],
	OPtr interfaces.BlobPtr[O],
] struct {
	awf    interfaces.BlobWriterFactory
	parser Parser[O, OPtr]
}

func MakeReaderAkteParseSaver[
	O interfaces.Blob[O],
	OPtr interfaces.BlobPtr[O],
](
	awf interfaces.BlobWriterFactory,
	parser Parser[O, OPtr],
) readerAkteParseSaver[O, OPtr] {
	return readerAkteParseSaver[O, OPtr]{
		awf:    awf,
		parser: parser,
	}
}

func (f readerAkteParseSaver[O, OPtr]) ParseSaveAkte(
	r io.Reader,
	t OPtr,
) (sh interfaces.ShaLike, n int64, err error) {
	var (
		aw  sha.WriteCloser
		sh1 interfaces.ShaLike
	)

	if aw, err = f.awf.BlobWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, aw)

	if n, err = f.ParseAkte(io.TeeReader(r, aw), t); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh = aw.GetShaLike()

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

func (f readerAkteParseSaver[O, OPtr]) ParseAkte(
	r io.Reader,
	t OPtr,
) (n int64, err error) {
	var n1 int64

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

		if n1, err = f.parser.ParseBlob(pr, t); err != nil {
			pr.CloseWithError(err)
		}
	}(pr)

	if n, err = io.Copy(pw, r); err != nil {
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

	return
}
