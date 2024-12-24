package blob_store

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
)

type savedBlobFormatter struct {
	arf interfaces.BlobReaderFactory
}

func MakeSavedBlobFormatter(
	blobReaderFactory interfaces.BlobReaderFactory,
) savedBlobFormatter {
	return savedBlobFormatter{
		arf: blobReaderFactory,
	}
}

func (f savedBlobFormatter) FormatSavedBlob(
	w io.Writer,
	sh interfaces.Sha,
) (n int64, err error) {
	var ar sha.ReadCloser

	if ar, err = f.arf.BlobReader(sh); err != nil {
		if errors.IsNotExist(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	defer errors.DeferredCloser(&err, ar)

	if n, err = io.Copy(w, ar); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
