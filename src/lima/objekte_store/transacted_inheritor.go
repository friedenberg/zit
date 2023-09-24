package objekte_store

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/pool"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type TransactedInheritor interface {
	InflateFromDataIdentityAndStoreAndInherit(*sku.Transacted) error
}

type InflatorStorer interface {
	TransactedDataIdentityInflator
	AkteStorer[*sku.Transacted]
}

type Inheritor interface {
	Inherit(*sku.Transacted) error
}

type heritableElement interface{}

type heritableElementPtr[T any] interface {
	schnittstellen.PoolablePtr[T]
}

type transactedInheritor[T heritableElement, TPtr heritableElementPtr[T]] struct {
	inflatorStorer InflatorStorer
	inheritor      Inheritor
}

func MakeTransactedInheritor[T heritableElement, TPtr heritableElementPtr[T]](
	inflatorStorer InflatorStorer,
	inheritor Inheritor,
) *transactedInheritor[T, TPtr] {
	return &transactedInheritor[T, TPtr]{
		inflatorStorer: inflatorStorer,
		inheritor:      inheritor,
	}
}

func (ti *transactedInheritor[T, TPtr]) InflateFromDataIdentityAndStoreAndInherit(
	sk *sku.Transacted,
) (err error) {
	var t *sku.Transacted

	if t, err = ti.inflatorStorer.InflateFromSku(sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	// if err = ti.inflatorStorer.StoreAkte(t); err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	shouldRepool := true

	if err = ti.inheritor.Inherit(t); err != nil {
		if pool.IsDoNotRepool(err) {
			shouldRepool = false
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	if shouldRepool {
		sku.GetTransactedPool().Put(t)
	}

	return
}
