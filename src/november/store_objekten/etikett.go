package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/kilo/objekte_store"
	"github.com/friedenberg/zit/src/mike/store_util"
)

type etikettStore struct {
	store_util.StoreUtil
}

func (s etikettStore) ReadOne(
	k schnittstellen.StringerGattungGetter,
) (tt *sku.Transacted, err error) {
	var e kennung.Etikett

	if err = e.Set(k.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	tt1, ok := s.StoreUtil.GetKonfig().GetEtikett(e)

	if !ok {
		err = errors.Wrap(objekte_store.ErrNotFound{Id: e})
		return
	}

	tt = sku.GetTransactedPool().Get()

	if err = tt.SetFromSkuLike(tt1); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s etikettStore) ReadAllSchwanzen(
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	if err = s.StoreUtil.GetKonfig().EachEtikett(f); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s etikettStore) ReadAll(
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	eachSku := func(o *sku.Transacted) (err error) {
		if o.GetGattung() != gattung.Etikett {
			return
		}

		if err = f(o); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	if err = s.StoreUtil.GetBestandsaufnahmeStore().ReadAllSkus(eachSku); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
