package objekte_store

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/sha"
)

// TODO-P3 rename to AkteSaver
type AkteTextSaver[
	T schnittstellen.Objekte[T],
	T1 schnittstellen.ObjektePtr[T],
] interface {
	SaveAkteText(T1) (int64, error)
}

type akteTextSaver[
	T schnittstellen.Objekte[T],
	T1 schnittstellen.ObjektePtr[T],
] struct {
	awf           schnittstellen.AkteWriterFactory
	akteFormatter schnittstellen.Formatter[T, T1]
}

func MakeAkteTextSaver[
	T schnittstellen.Objekte[T],
	T1 schnittstellen.ObjektePtr[T],
](
	awf schnittstellen.AkteWriterFactory,
	akteFormatter schnittstellen.Formatter[T, T1],
) *akteTextSaver[T, T1] {
	return &akteTextSaver[T, T1]{
		awf:           awf,
		akteFormatter: akteFormatter,
	}
}

func (h *akteTextSaver[T, T1]) SaveAkteText(
	o T1,
) (n int64, err error) {
	var w sha.WriteCloser

	if w, err = h.awf.AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, w.Close)

	if n, err = h.akteFormatter.Format(w, o); err != nil {
		err = errors.Wrap(err)
		return
	}

	o.SetAkteSha(w.Sha())

	return
}
