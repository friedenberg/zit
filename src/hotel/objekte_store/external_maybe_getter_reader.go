package objekte_store

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/hotel/objekte"
)

type ExternalMaybeGetterReader[
	O objekte.Akte[O],
	OPtr objekte.AktePtr[O],
	K kennung.KennungLike[K],
	KPtr kennung.KennungLikePtr[K],
] interface {
	ReadOne(
		sku.Transacted[K, KPtr],
	) (*objekte.CheckedOut[K, KPtr], error)
}

type externalMaybeGetterReader[
	O objekte.Akte[O],
	OPtr objekte.AktePtr[O],
	K kennung.KennungLike[K],
	KPtr kennung.KennungLikePtr[K],
] struct {
	getter func(KPtr) (*sku.ExternalMaybe[K, KPtr], bool)
	ExternalReader[*sku.ExternalMaybe[K, KPtr],
		*sku.Transacted[K, KPtr],
		sku.External[K, KPtr]]
}

func MakeExternalMaybeGetterReader[
	O objekte.Akte[O],
	OPtr objekte.AktePtr[O],
	K kennung.KennungLike[K],
	KPtr kennung.KennungLikePtr[K],
](
	getter func(KPtr) (*sku.ExternalMaybe[K, KPtr], bool),
	er ExternalReader[
		*sku.ExternalMaybe[K, KPtr],
		*sku.Transacted[K, KPtr],
		sku.External[K, KPtr],
	],
) ExternalMaybeGetterReader[O, OPtr, K, KPtr] {
	return externalMaybeGetterReader[O, OPtr, K, KPtr]{
		getter:         getter,
		ExternalReader: er,
	}
}

func (emgr externalMaybeGetterReader[O, OPtr, K, KPtr]) ReadOne(
	i sku.Transacted[K, KPtr],
) (co *objekte.CheckedOut[K, KPtr], err error) {
	co = &objekte.CheckedOut[K, KPtr]{
		Internal: i,
	}

	ok := false

	var e *sku.ExternalMaybe[K, KPtr]

	if e, ok = emgr.getter(KPtr(i.GetKennungPtr())); !ok {
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
