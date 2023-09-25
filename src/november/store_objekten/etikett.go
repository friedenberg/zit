package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/etikett_akte"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/kilo/objekte_store"
	"github.com/friedenberg/zit/src/mike/store_util"
)

type EtikettTransactedReader = objekte_store.TransactedReader

type etikettStore struct {
	*store_util.CommonStore[
		etikett_akte.V0,
		*etikett_akte.V0,
	]
}

func makeEtikettStore(
	sa store_util.StoreUtil,
) (s *etikettStore, err error) {
	s = &etikettStore{}

	s.CommonStore, err = store_util.MakeCommonStore[
		etikett_akte.V0,
		*etikett_akte.V0,
	](
		gattung.Etikett,
		s,
		sa,
		s,
	)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	newOrUpdated := func(t *sku.Transacted) (err error) {
		s.StoreUtil.CommitUpdatedTransacted(t)

		if err = s.StoreUtil.GetKonfigPtr().AddEtikett(t); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	s.CommonStore.CreateOrUpdater = objekte_store.MakeCreateOrUpdate(
		sa,
		sa.GetStandort().GetLockSmith(),
		sa.GetStandort(),
		EtikettTransactedReader(s),
		objekte_store.CreateOrUpdateDelegate{
			New: func(t *sku.Transacted) (err error) {
				if err = newOrUpdated(t); err != nil {
					err = errors.Wrap(err)
					return
				}

				return s.LogWriter.New(t)
			},
			Updated: func(t *sku.Transacted) (err error) {
				if err = newOrUpdated(t); err != nil {
					err = errors.Wrap(err)
					return
				}

				return s.LogWriter.Updated(t)
			},
			Unchanged: func(t *sku.Transacted) (err error) {
				return s.LogWriter.Unchanged(t)
			},
		},
		sa.GetAbbrStore(),
		sa.GetPersistentMetadateiFormat(),
		objekte_format.Options{IncludeTai: true},
		sa,
	)

	return
}

func (s etikettStore) Flush() (err error) {
	return
}

func (s etikettStore) AddOne(t *sku.Transacted) (err error) {
	s.StoreUtil.GetKonfigPtr().AddEtikett(t)
	return
}

func (s etikettStore) UpdateOne(t *sku.Transacted) (err error) {
	s.StoreUtil.GetKonfigPtr().AddEtikett(t)
	return
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
