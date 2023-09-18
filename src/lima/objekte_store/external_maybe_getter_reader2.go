package objekte_store

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/juliett/objekte"
)

type ExternalMaybeGetterReader2 interface {
	ReadOne(*sku.Transacted2) (*objekte.CheckedOut2, error)
}

type externalMaybeGetterReader2 struct {
	getter func(kennung.Kennung2) (*sku.ExternalMaybe, bool)
	ExternalReader[
		*sku.ExternalMaybe,
		*sku.Transacted2,
		*sku.External2,
	]
}

func MakeExternalMaybeGetterReader2(
	getter func(kennung.Kennung2) (*sku.ExternalMaybe, bool),
	er ExternalReader[
		*sku.ExternalMaybe,
		*sku.Transacted2,
		*sku.External2,
	],
) ExternalMaybeGetterReader2 {
	return externalMaybeGetterReader2{
		getter:         getter,
		ExternalReader: er,
	}
}

func (emgr externalMaybeGetterReader2) ReadOne(
	sk2 *sku.Transacted2,
) (co *objekte.CheckedOut2, err error) {
	co = &objekte.CheckedOut2{
		Internal: *sk2,
	}

	ok := false

	var e *sku.ExternalMaybe

	if e, ok = emgr.getter(sk2.Kennung); !ok {
		err = iter.MakeErrStopIteration()
		return
	}

	var e2 *sku.External2

	if e2, err = emgr.ReadOneExternal(e, sk2); err != nil {
		if errors.IsNotExist(err) {
			err = iter.MakeErrStopIteration()
		} else {
			err = errors.Wrapf(err, "Cwd: %#v", e)
		}

		return
	}

	co.External = *e2

	return
}
