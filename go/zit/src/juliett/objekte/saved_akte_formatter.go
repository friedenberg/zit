package objekte

import (
	"io"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/delta/sha"
)

type savedAkteFormatter struct {
	arf schnittstellen.AkteReaderFactory
}

func MakeSavedAkteFormatter(
	akteReaderFactory schnittstellen.AkteReaderFactory,
) savedAkteFormatter {
	return savedAkteFormatter{
		arf: akteReaderFactory,
	}
}

func (f savedAkteFormatter) FormatSavedAkte(
	w io.Writer,
	sh schnittstellen.ShaLike,
) (n int64, err error) {
	var ar sha.ReadCloser

	if ar, err = f.arf.AkteReader(sh); err != nil {
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
