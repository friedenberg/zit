package objekte

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/golf/sku"
)

type TransactedInheritor interface {
	InflateFromDataIdentityAndStoreAndInherit(sku.DataIdentity) error
}

type InflatorStorer[T any] interface {
	TransactedDataIdentityInflator[T]
	ObjekteStorer[T]
	AkteStorer[T]
}

type Inheritor[T any] interface {
	Inherit(T) error
}

type heritableElement interface {
	gattung.Element
}

type heritableElementPtr[T gattung.Element] interface {
	gattung.ElementPtr[T]
}

type transactedInheritor[T heritableElement, TPtr heritableElementPtr[T]] struct {
	inflatorStorer InflatorStorer[TPtr]
	inheritor      Inheritor[TPtr]
	pool           collections.Pool2Like[T, TPtr]
}

func MakeTransactedInheritor[T heritableElement, TPtr heritableElementPtr[T]](
	inflatorStorer InflatorStorer[TPtr],
	inheritor Inheritor[TPtr],
	pool collections.Pool2Like[T, TPtr],
) *transactedInheritor[T, TPtr] {
	return &transactedInheritor[T, TPtr]{
		inflatorStorer: inflatorStorer,
		inheritor:      inheritor,
		pool:           pool,
	}
}

func (ti *transactedInheritor[T, TPtr]) InflateFromDataIdentityAndStoreAndInherit(
	sk sku.DataIdentity,
) (err error) {
	var t TPtr

	if t, err = ti.inflatorStorer.InflateFromDataIdentity(sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	// if err = ti.inflatorStorer.StoreAkte(t); err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	if err = ti.inflatorStorer.StoreObjekte(t); err != nil {
		err = errors.Wrap(err)
		return
	}

	shouldRepool := true

	if err = ti.inheritor.Inherit(t); err != nil {
		if collections.IsDoNotRepool(err) {
			shouldRepool = false
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	if shouldRepool {
		ti.pool.Put(t)
	}

	return
}
