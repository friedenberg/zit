package objekte

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/echo/sha"
)

// TODO-P3 rename to AkteSaver
type AkteTextSaver[
	T gattung.Objekte[T],
	T1 gattung.ObjektePtr[T],
] interface {
	SaveAkteText(T1) (int64, error)
}

type akteTextSaver[
	T gattung.Objekte[T],
	T1 gattung.ObjektePtr[T],
] struct {
	awf           gattung.AkteWriterFactory
	akteFormatter gattung.Formatter[T, T1]
}

func MakeAkteTextSaver[
	T gattung.Objekte[T],
	T1 gattung.ObjektePtr[T],
](
	awf gattung.AkteWriterFactory,
	akteFormatter gattung.Formatter[T, T1],
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
