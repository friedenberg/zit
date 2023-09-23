package objekte_store

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/juliett/objekte"
)

type ExternalMaybeGetterReader[
	O objekte.Akte[O],
	OPtr objekte.AktePtr[O],
	K kennung.KennungLike[K],
	KPtr kennung.KennungLikePtr[K],
] interface {
	ReadOne(
		sku.Transacted,
	) (*objekte.CheckedOut2, error)
}

type externalMaybeGetterReader[
	O objekte.Akte[O],
	OPtr objekte.AktePtr[O],
	K kennung.KennungLike[K],
	KPtr kennung.KennungLikePtr[K],
] struct {
	getter func(kennung.Kennung) (*sku.ExternalMaybe, bool)
	ExternalReader[
		sku.ExternalMaybe,
		sku.SkuLikePtr,
	]
}

func MakeExternalMaybeGetterReader[
	O objekte.Akte[O],
	OPtr objekte.AktePtr[O],
	K kennung.KennungLike[K],
	KPtr kennung.KennungLikePtr[K],
](
	getter func(kennung.Kennung) (*sku.ExternalMaybe, bool),
	er ExternalReader[
		sku.ExternalMaybe,
		sku.SkuLikePtr,
	],
) ExternalMaybeGetterReader[O, OPtr, K, KPtr] {
	return externalMaybeGetterReader[O, OPtr, K, KPtr]{
		getter:         getter,
		ExternalReader: er,
	}
}

func (emgr externalMaybeGetterReader[O, OPtr, K, KPtr]) ReadOne(
	i sku.Transacted,
) (co *objekte.CheckedOut2, err error) {
	co = &objekte.CheckedOut2{
		Internal: i,
	}

	ok := false

	var e *sku.ExternalMaybe

	if e, ok = emgr.getter(i.GetKennungPtr()); !ok {
		err = iter.MakeErrStopIteration()
		return
	}

	var ex *sku.External

	if ex, err = emgr.ReadOneExternal(*e, &i); err != nil {
		if errors.IsNotExist(err) {
			err = iter.MakeErrStopIteration()
		} else {
			err = errors.Wrapf(err, "Cwd: %#v", e)
		}

		return
	}

	co.External = *ex

	return
}
