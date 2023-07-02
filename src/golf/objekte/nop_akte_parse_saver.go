package objekte

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/sha"
)

type nopAkteParseSaver[
	O Akte[O],
	OPtr AktePtr[O],
] struct {
	awf schnittstellen.AkteWriterFactory
}

func MakeNopAkteParseSaver[
	O Akte[O],
	OPtr AktePtr[O],
](awf schnittstellen.AkteWriterFactory,
) nopAkteParseSaver[O, OPtr] {
	return nopAkteParseSaver[O, OPtr]{
		awf: awf,
	}
}

func (f nopAkteParseSaver[O, OPtr]) ParseSaveAkte(
	r io.Reader,
	t OPtr,
) (sh schnittstellen.ShaLike, n int64, err error) {
	var aw sha.WriteCloser

	if aw, err = f.awf.AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, aw)

	if n, err = io.Copy(aw, r); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh = aw.GetShaLike()

	return
}
