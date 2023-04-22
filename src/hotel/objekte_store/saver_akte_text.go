package objekte_store

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/golf/objekte"
)

// TODO-P3 rename to AkteSaver
type AkteTextSaver[
	T objekte.Objekte[T],
	T1 objekte.ObjektePtr[T],
] interface {
	SaveAkteText(T) (schnittstellen.Sha, int64, error)
}

type akteTextSaver[
	T objekte.Objekte[T],
	T1 objekte.ObjektePtr[T],
] struct {
	awf        schnittstellen.AkteWriterFactory
	akteFormat objekte.AkteFormat[T, T1]
}

func MakeAkteTextSaver[
	T objekte.Objekte[T],
	T1 objekte.ObjektePtr[T],
](
	awf schnittstellen.AkteWriterFactory,
	akteFormat objekte.AkteFormat[T, T1],
) akteTextSaver[T, T1] {
	return akteTextSaver[T, T1]{
		awf:        awf,
		akteFormat: akteFormat,
	}
}

func (h akteTextSaver[T, T1]) SaveAkteText(
	o T,
) (sh schnittstellen.Sha, n int64, err error) {
	var w sha.WriteCloser

	if w, err = h.awf.AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, w)

	if n, err = h.akteFormat.FormatParsedAkte(w, o); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh = sha.Make(w.Sha())

	return
}
