package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/alfa/toml"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/golf/transaktion"
	"github.com/friedenberg/zit/src/hotel/kasten"
	"github.com/friedenberg/zit/src/hotel/objekte"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
	"github.com/friedenberg/zit/src/kilo/store_util"
)

type KastenStore interface {
	CommonStore[
		kasten.Akte,
		*kasten.Akte,
		kennung.Kasten,
		*kennung.Kasten,
	]
}

type KastenTransactedReader = objekte_store.TransactedReader[
	*kennung.Kasten,
	*kasten.Transacted,
]

type kastenStore struct {
	*commonStore[
		kasten.Akte,
		*kasten.Akte,
		kennung.Kasten,
		*kennung.Kasten,
	]
}

func makeKastenStore(
	sa store_util.StoreUtil,
) (s *kastenStore, err error) {
	s = &kastenStore{}

	s.commonStore, err = makeCommonStore[
		kasten.Akte,
		*kasten.Akte,
		kennung.Kasten,
		*kennung.Kasten,
	](
		gattung.Kasten,
		s,
		sa,
		s,
		objekte_store.MakeAkteFormat[kasten.Akte, *kasten.Akte](
			objekte.MakeTextParserIgnoreTomlErrors[kasten.Akte](sa),
			objekte.ParsedAkteTomlFormatter[kasten.Akte]{},
			sa,
		),
	)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	newOrUpdated := func(t *kasten.Transacted) (err error) {
		s.StoreUtil.CommitUpdatedTransacted(t)

		if err = s.StoreUtil.GetKonfigPtr().AddKasten(t); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	s.commonStore.CreateOrUpdater = objekte_store.MakeCreateOrUpdate[
		kasten.Akte,
		*kasten.Akte,
		kennung.Kasten,
		*kennung.Kasten,
	](
		sa,
		sa.GetLockSmith(),
		sa.ObjekteReaderWriterFactory(gattung.Kasten),
		sa,
		KastenTransactedReader(s),
		objekte_store.CreateOrUpdateDelegate[*kasten.Transacted]{
			New: func(t *kasten.Transacted) (err error) {
				if err = newOrUpdated(t); err != nil {
					err = errors.Wrap(err)
					return
				}

				return s.LogWriter.New(t)
			},
			Updated: func(t *kasten.Transacted) (err error) {
				if err = newOrUpdated(t); err != nil {
					err = errors.Wrap(err)
					return
				}

				return s.LogWriter.Updated(t)
			},
			Unchanged: func(t *kasten.Transacted) (err error) {
				return s.LogWriter.Unchanged(t)
			},
		},
		sa.GetAbbrStore(),
		sa.GetPersistentMetadateiFormat(),
		sa,
	)

	return
}

func (s kastenStore) Flush() (err error) {
	return
}

func (s kastenStore) addOne(t *kasten.Transacted) (err error) {
	s.StoreUtil.GetKonfigPtr().AddKasten(t)
	return
}

func (s kastenStore) updateOne(t *kasten.Transacted) (err error) {
	s.StoreUtil.GetKonfigPtr().AddKasten(t)
	return
}

// TODO-P3
func (s kastenStore) ReadAllSchwanzen(
	f schnittstellen.FuncIter[*kasten.Transacted],
) (err error) {
	// TODO-P2 switch to pointers
	if err = s.StoreUtil.GetKonfig().Kisten.Each(
		func(e kasten.Transacted) (err error) {
			return f(&e)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s kastenStore) ReadAll(
	f schnittstellen.FuncIter[*kasten.Transacted],
) (err error) {
	eachSku := func(sk sku.SkuLikePtr) (err error) {
		if sk.GetGattung() != gattung.Kasten {
			return
		}

		var te *kasten.Transacted

		if te, err = s.InflateFromSku(sk); err != nil {
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

func (s kastenStore) ReadOne(
	k *kennung.Kasten,
) (tt *kasten.Transacted, err error) {
	errors.TodoP3("add support for working directory")
	errors.TodoP3("inherited-kastenen-etiketten")
	tt = s.StoreUtil.GetKonfig().GetKasten(*k)

	if tt == nil {
		err = errors.Wrap(objekte_store.ErrNotFound{Id: k})
		return
	}

	return
}
