package akten

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
)

type nopAkteParseSaver[
	O schnittstellen.Akte[O],
	OPtr schnittstellen.AktePtr[O],
] struct {
	awf schnittstellen.AkteWriterFactory
}

func MakeNopAkteParseSaver[
	O schnittstellen.Akte[O],
	OPtr schnittstellen.AktePtr[O],
](awf schnittstellen.AkteWriterFactory,
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

	if aw, err = f.awf.AkteWriter(); err != nil {
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
