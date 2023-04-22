package objekte_store

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/golf/persisted_metadatei_format"
)

type ObjekteSaver[
	T objekte.Objekte[T],
	T1 objekte.ObjektePtr[T],
] interface {
	SaveObjekte(T1) (schnittstellen.Sha, error)
}

type objekteSaver[
	T objekte.Objekte[T],
	T1 objekte.ObjektePtr[T],
] struct {
	writerFactory schnittstellen.ObjekteWriterFactory
	formatter     persisted_metadatei_format.Format
}

func MakeObjekteSaver[
	T objekte.Objekte[T],
	T1 objekte.ObjektePtr[T],
](
	writerFactory schnittstellen.ObjekteWriterFactory,
	pmf persisted_metadatei_format.Format,
) *objekteSaver[T, T1] {
	if writerFactory == nil {
		panic("schnittstellen.ObjekteWriterFactory was nil")
	}

	if pmf == nil {
		panic("persisted_metadatei_format.Format was nil")
	}

	return &objekteSaver[T, T1]{
		writerFactory: writerFactory,
		formatter:     pmf,
	}
}

func (h *objekteSaver[T, T1]) SaveObjekte(
	o T1,
) (sh schnittstellen.Sha, err error) {
	var w sha.WriteCloser

	if w, err = h.writerFactory.ObjekteWriter(); err != nil {
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
