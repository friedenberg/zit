package typed_blob_store

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
)

type nopBlobParseSaver[
	O any,
	OPtr interfaces.Ptr[O],
] struct {
	awf interfaces.BlobWriter
}

func MakeNopBlobParseSaver[
	O any,
	OPtr interfaces.Ptr[O],
](awf interfaces.BlobWriter,
) nopBlobParseSaver[O, OPtr] {
	return nopBlobParseSaver[O, OPtr]{
		awf: awf,
	}
}

func (f nopBlobParseSaver[O, OPtr]) ParseBlob(
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
