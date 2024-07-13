package akten

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
)

type savedAkteFormatter struct {
	arf interfaces.BlobReaderFactory
}

func MakeSavedAkteFormatter(
	akteReaderFactory interfaces.BlobReaderFactory,
) savedAkteFormatter {
	return savedAkteFormatter{
		arf: akteReaderFactory,
	}
}

func (f savedAkteFormatter) FormatSavedAkte(
	w io.Writer,
	sh interfaces.ShaLike,
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
