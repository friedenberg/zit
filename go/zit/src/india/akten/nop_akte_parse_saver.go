package akten

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
)

type nopAkteParseSaver[
	O interfaces.Blob[O],
	OPtr interfaces.BlobPtr[O],
] struct {
	awf interfaces.BlobWriterFactory
}

func MakeNopAkteParseSaver[
	O interfaces.Blob[O],
	OPtr interfaces.BlobPtr[O],
](awf interfaces.BlobWriterFactory,
) nopAkteParseSaver[O, OPtr] {
	return nopAkteParseSaver[O, OPtr]{
		awf: awf,
	}
}

func (f nopAkteParseSaver[O, OPtr]) ParseAkte(
	r io.Reader,
	t OPtr,
) (n int64, err error) {
	var aw sha.WriteCloser

	if aw, err = f.awf.BlobWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, aw)

	if n, err = io.Copy(aw, r); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
