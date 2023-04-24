package objekte_store

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/objekte"
)

type ExternalMaybeGetterReader[
	O objekte.Objekte[O],
	OPtr objekte.ObjektePtr[O],
	K schnittstellen.Id[K],
	KPtr schnittstellen.IdPtr[K],
	V any,
	VPtr objekte.VerzeichnissePtr[V, O],
] interface {
	ReadOne(
		objekte.Transacted[O, OPtr, K, KPtr, V, VPtr],
	) (*objekte.CheckedOut[O, OPtr, K, KPtr, V, VPtr], error)
}

type externalMaybeGetterReader[
	O objekte.Objekte[O],
	OPtr objekte.ObjektePtr[O],
	K schnittstellen.Id[K],
	KPtr schnittstellen.IdPtr[K],
	V any,
	VPtr objekte.VerzeichnissePtr[V, O],
] struct {
	getter func(K) (sku.ExternalMaybe[K, KPtr], bool)
	ExternalReader[sku.ExternalMaybe[K, KPtr],
		*objekte.Transacted[O, OPtr, K, KPtr, V, VPtr],
		objekte.External[O, OPtr, K, KPtr]]
}

func MakeExternalMaybeGetterReader[
	O objekte.Objekte[O],
	OPtr objekte.ObjektePtr[O],
	K schnittstellen.Id[K],
	KPtr schnittstellen.IdPtr[K],
	V any,
	VPtr objekte.VerzeichnissePtr[V, O],
](
	getter func(K) (sku.ExternalMaybe[K, KPtr], bool),
	er ExternalReader[
		sku.ExternalMaybe[K, KPtr],
		*objekte.Transacted[O, OPtr, K, KPtr, V, VPtr],
		objekte.External[O, OPtr, K, KPtr],
	],
) ExternalMaybeGetterReader[O, OPtr, K, KPtr, V, VPtr] {
	return externalMaybeGetterReader[
		O,
		OPtr,
		K,
		KPtr,
		V,
		VPtr,
	]{
		getter:         getter,
		ExternalReader: er,
	}
}

func (emgr externalMaybeGetterReader[O, OPtr, K, KPtr, V, VPtr]) ReadOne(
	i objekte.Transacted[O, OPtr, K, KPtr, V, VPtr],
) (co *objekte.CheckedOut[O, OPtr, K, KPtr, V, VPtr], err error) {
	co = &objekte.CheckedOut[O, OPtr, K, KPtr, V, VPtr]{
		Internal: i,
	}

	ok := false

	var e sku.ExternalMaybe[K, KPtr]

	if e, ok = emgr.getter(i.Sku.Kennung); !ok {
		err = iter.MakeErrStopIteration()
		return
	}

	if co.External, err = emgr.ReadOneExternal(e, &i); err != nil {
		if errors.IsNotExist(err) {
			err = iter.MakeErrStopIteration()
		} else {
			err = errors.Wrapf(err, "Cwd: %#v", e)
		}

		return
	}

	return
}
