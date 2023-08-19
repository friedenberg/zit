package objekte

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/sha"
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
