package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/alfa/toml"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/log"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/golf/transaktion"
	"github.com/friedenberg/zit/src/hotel/objekte"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
	"github.com/friedenberg/zit/src/hotel/typ"
	"github.com/friedenberg/zit/src/kilo/store_util"
)

type TypStore interface {
	CommonStore[
		typ.Akte,
		*typ.Akte,
		kennung.Typ,
		*kennung.Typ,
	]
}

type TypTransactedReader = objekte_store.TransactedReader[
	*kennung.Typ,
	*typ.Transacted,
]

type typStore struct {
	*commonStore[
		typ.Akte,
		*typ.Akte,
		kennung.Typ,
		*kennung.Typ,
	]
}

func makeTypStore(
	sa store_util.StoreUtil,
) (s *typStore, err error) {
	s = &typStore{}

	s.commonStore, err = makeCommonStore[
		typ.Akte,
		*typ.Akte,
		kennung.Typ,
		*kennung.Typ,
	](
		gattung.Typ,
		s,
		sa,
		s,
		objekte_store.MakeAkteFormat[typ.Akte, *typ.Akte](
			objekte.MakeTextParserIgnoreTomlErrors[typ.Akte](sa),
			objekte.ParsedAkteTomlFormatter[typ.Akte]{},
			sa,
		),
	)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	newOrUpdated := func(t *typ.Transacted) (err error) {
		s.StoreUtil.CommitUpdatedTransacted(t)

		if err = s.StoreUtil.GetKonfigPtr().AddTyp(t); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	s.commonStore.CreateOrUpdater = objekte_store.MakeCreateOrUpdate[
		typ.Akte,
		*typ.Akte,
		kennung.Typ,
		*kennung.Typ,
	](
		sa,
		sa.GetLockSmith(),
		s.commonStore,
		sa,
		TypTransactedReader(s),
		objekte_store.CreateOrUpdateDelegate[*typ.Transacted]{
			New: func(t *typ.Transacted) (err error) {
				if err = newOrUpdated(t); err != nil {
					err = errors.Wrap(err)
					return
				}

				return s.LogWriter.New(t)
			},
			Updated: func(t *typ.Transacted) (err error) {
				if err = newOrUpdated(t); err != nil {
					err = errors.Wrap(err)
					return
				}

				return s.LogWriter.Updated(t)
			},
			Unchanged: func(t *typ.Transacted) (err error) {
				return s.LogWriter.Unchanged(t)
			},
		},
		sa.GetAbbrStore(),
		sa.GetPersistentMetadateiFormat(),
		sa,
	)

	return
}

func (s typStore) Flush() (err error) {
	return
}

func (s typStore) addOne(t *typ.Transacted) (err error) {
	s.StoreUtil.GetKonfigPtr().AddTyp(t)
	return
}

func (s typStore) updateOne(t *typ.Transacted) (err error) {
	log.Log().Printf("adding one: %s", t.GetSkuLike())
	s.StoreUtil.GetKonfigPtr().AddTyp(t)
	log.Log().Printf("done adding one: %s", t.GetSkuLike())
	return
}

// TODO-P3
func (s typStore) ReadAllSchwanzen(
	f schnittstellen.FuncIter[*typ.Transacted],
) (err error) {
	// TODO-P2 switch to pointers
	if err = s.StoreUtil.GetKonfig().Typen.Each(
		func(e typ.Transacted) (err error) {
			return f(&e)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s typStore) ReadAll(
	f schnittstellen.FuncIter[*typ.Transacted],
) (err error) {
	eachSku := func(sk sku.SkuLikePtr) (err error) {
		if sk.GetGattung() != gattung.Typ {
			return
		}

		var te *typ.Transacted

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

func (s typStore) ReadOne(
	k *kennung.Typ,
) (tt *typ.Transacted, err error) {
	errors.TodoP3("add support for working directory")
	errors.TodoP3("inherited-typen-etiketten")
	log.Log().Printf("reading: %s", k)
	tt = s.StoreUtil.GetKonfig().GetApproximatedTyp(*k).ActualOrNil()

	if tt == nil {
		err = errors.Wrap(objekte_store.ErrNotFound{Id: k})
		return
	}

	return
}
