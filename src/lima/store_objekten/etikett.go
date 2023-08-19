package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/alfa/toml"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/golf/transaktion"
	"github.com/friedenberg/zit/src/hotel/etikett"
	"github.com/friedenberg/zit/src/hotel/objekte"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
	"github.com/friedenberg/zit/src/kilo/store_util"
)

type EtikettStore interface {
	CommonStore[
		etikett.Akte,
		*etikett.Akte,
		kennung.Etikett,
		*kennung.Etikett,
	]
}

type EtikettTransactedReader = objekte_store.TransactedReader[
	*kennung.Etikett,
	*sku.TransactedEtikett,
]

type etikettStore struct {
	*commonStore[
		etikett.Akte,
		*etikett.Akte,
		kennung.Etikett,
		*kennung.Etikett,
	]
}

func makeEtikettStore(
	sa store_util.StoreUtil,
) (s *etikettStore, err error) {
	s = &etikettStore{}

	s.commonStore, err = makeCommonStore[
		etikett.Akte,
		*etikett.Akte,
		kennung.Etikett,
		*kennung.Etikett,
	](
		gattung.Etikett,
		s,
		sa,
		s,
		objekte_store.MakeAkteFormat[etikett.Akte, *etikett.Akte](
			objekte.MakeTextParserIgnoreTomlErrors[etikett.Akte](sa),
			objekte.ParsedAkteTomlFormatter[etikett.Akte]{},
			sa,
		),
	)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	newOrUpdated := func(t *sku.TransactedEtikett) (err error) {
		s.StoreUtil.CommitUpdatedTransacted(t)

		if err = s.StoreUtil.GetKonfigPtr().AddEtikett(t); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	s.commonStore.CreateOrUpdater = objekte_store.MakeCreateOrUpdate[
		etikett.Akte,
		*etikett.Akte,
		kennung.Etikett,
		*kennung.Etikett,
	](
		sa,
		sa.GetLockSmith(),
		sa.ObjekteReaderWriterFactory(gattung.Etikett),
		sa,
		EtikettTransactedReader(s),
		objekte_store.CreateOrUpdateDelegate[*sku.TransactedEtikett]{
			New: func(t *sku.TransactedEtikett) (err error) {
				if err = newOrUpdated(t); err != nil {
					err = errors.Wrap(err)
					return
				}

				return s.LogWriter.New(t)
			},
			Updated: func(t *sku.TransactedEtikett) (err error) {
				if err = newOrUpdated(t); err != nil {
					err = errors.Wrap(err)
					return
				}

				return s.LogWriter.Updated(t)
			},
			Unchanged: func(t *sku.TransactedEtikett) (err error) {
				return s.LogWriter.Unchanged(t)
			},
		},
		sa.GetAbbrStore(),
		sa.GetPersistentMetadateiFormat(),
		sa,
	)

	return
}

func (s etikettStore) Flush() (err error) {
	return
}

func (s etikettStore) addOne(t *sku.TransactedEtikett) (err error) {
	s.StoreUtil.GetKonfigPtr().AddEtikett(t)
	return
}

func (s etikettStore) updateOne(t *sku.TransactedEtikett) (err error) {
	s.StoreUtil.GetKonfigPtr().AddEtikett(t)
	return
}

func (s etikettStore) ReadOne(
	k *kennung.Etikett,
) (tt *sku.TransactedEtikett, err error) {
	tt1 := s.StoreUtil.GetKonfig().GetEtikett(*k)
	tt = &tt1

	if tt == nil {
		err = errors.Wrap(objekte_store.ErrNotFound{Id: k})
		return
	}

	return
}

func (s etikettStore) ReadAllSchwanzen(
	f schnittstellen.FuncIter[*sku.TransactedEtikett],
) (err error) {
	if err = s.StoreUtil.GetKonfig().EachEtikett(f); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s etikettStore) ReadAll(
	f schnittstellen.FuncIter[*sku.TransactedEtikett],
) (err error) {
	eachSku := func(o sku.SkuLikePtr) (err error) {
		if o.GetGattung() != gattung.Etikett {
			return
		}

		var te *sku.TransactedEtikett

		if te, err = s.InflateFromSku(o); err != nil {
			if errors.Is(err, toml.Error{}) {
				err = nil
			} else {
				err = errors.Wrap(err)
				return
			}
		}

		if err = f(te); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	if s.StoreUtil.GetKonfig().UseBestandsaufnahme {
		if err = s.StoreUtil.GetBestandsaufnahmeStore().ReadAllSkus(eachSku); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = s.StoreUtil.GetTransaktionStore().ReadAllTransaktions(
		func(t *transaktion.Transaktion) (err error) {
			if err = t.Skus.Each(
				eachSku,
			); err != nil {
				err = errors.Wrapf(
					err,
					"Transaktion: %s/%s: %s",
					t.Time.Kopf(),
					t.Time.Schwanz(),
					t.Time,
				)

				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
