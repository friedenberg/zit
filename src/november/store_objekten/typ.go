package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/log"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/kilo/objekte_store"
	"github.com/friedenberg/zit/src/mike/store_util"
)

type typStore struct {
	store_util.StoreUtil
}

func (s typStore) ReadOne(
	k1 schnittstellen.StringerGattungGetter,
) (tt *sku.Transacted, err error) {
	var k kennung.Typ

	if err = k.Set(k1.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.TodoP3("add support for working directory")
	errors.TodoP3("inherited-typen-etiketten")
	log.Log().Printf("reading: %s", k)
	t1 := s.StoreUtil.GetKonfig().GetApproximatedTyp(k).ActualOrNil()

	if t1 == nil {
		err = errors.Wrap(objekte_store.ErrNotFound{Id: k})
		return
	}

	tt = sku.GetTransactedPool().Get()

	if err = tt.SetFromSkuLike(t1); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
