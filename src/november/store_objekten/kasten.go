package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/kilo/objekte_store"
	"github.com/friedenberg/zit/src/mike/store_util"
)

type kastenStore struct {
	store_util.StoreUtil
}

func (s kastenStore) ReadOne(
	k schnittstellen.StringerGattungGetter,
) (tt *sku.Transacted, err error) {
	errors.TodoP3("add support for working directory")
	errors.TodoP3("inherited-kastenen-etiketten")
	var k1 kennung.Kasten

	if err = k1.Set(k.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	tt1 := s.StoreUtil.GetKonfig().GetKasten(k1)

	if tt1 == nil {
		err = errors.Wrap(objekte_store.ErrNotFound{Id: k1})
		return
	}

	tt = sku.GetTransactedPool().Get()

	if err = tt.SetFromSkuLike(tt1); err != nil {
		err = errors.Wrap(err)
	}

	return
}
