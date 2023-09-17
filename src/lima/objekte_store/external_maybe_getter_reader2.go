package objekte_store

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/juliett/objekte"
)

type ExternalMaybeGetterReader2 interface {
	ReadOne(sku.SkuLikePtr) (*objekte.CheckedOut2, error)
}

type externalMaybeGetterReader2 struct {
	getter func(kennung.Kennung) (*sku.ExternalMaybe, bool)
	ExternalReader[
		*sku.ExternalMaybe,
		sku.SkuLikePtr,
		*sku.External2,
	]
}

func MakeExternalMaybeGetterReader2(
	getter func(kennung.Kennung) (*sku.ExternalMaybe, bool),
	er ExternalReader[
		*sku.ExternalMaybe,
		sku.SkuLikePtr,
		*sku.External2,
	],
) ExternalMaybeGetterReader2 {
	return externalMaybeGetterReader2{
		getter:         getter,
		ExternalReader: er,
	}
}

func (emgr externalMaybeGetterReader2) ReadOne(
	i sku.SkuLikePtr,
) (co *objekte.CheckedOut2, err error) {
	co = &objekte.CheckedOut2{
		Internal: i,
	}

	ok := false

	var e *sku.ExternalMaybe

	if e, ok = emgr.getter(i.GetKennungLike()); !ok {
		err = iter.MakeErrStopIteration()
		return
	}

	if co.External, err = emgr.ReadOneExternal(e, i); err != nil {
		if errors.IsNotExist(err) {
			err = iter.MakeErrStopIteration()
		} else {
			err = errors.Wrapf(err, "Cwd: %#v", e)
		}

		return
	}

	return
}
