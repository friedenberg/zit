package objekte

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/schnittstellen"
)

type ObjekteSaver[
	T gattung.Objekte[T],
	T1 gattung.ObjektePtr[T],
] interface {
	SaveObjekte(T1) (schnittstellen.Sha, error)
}

type objekteSaver[
	T gattung.Objekte[T],
	T1 gattung.ObjektePtr[T],
] struct {
	writerFactory schnittstellen.ObjekteWriterFactory
	formatter     schnittstellen.Formatter[T, T1]
}

func MakeObjekteSaver[
	T gattung.Objekte[T],
	T1 gattung.ObjektePtr[T],
](
	writerFactory schnittstellen.ObjekteWriterFactory,
	formatter schnittstellen.Formatter[T, T1],
) *objekteSaver[T, T1] {
	return &objekteSaver[T, T1]{
		writerFactory: writerFactory,
		formatter:     formatter,
	}
}

func (h *objekteSaver[T, T1]) SaveObjekte(
	o T1,
) (sh schnittstellen.Sha, err error) {
	var w sha.WriteCloser

	if w, err = h.writerFactory.ObjekteWriter(
		o.GetGattung(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, w.Close)

	if _, err = h.formatter.Format(w, o); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh = w.Sha()

	return
}
