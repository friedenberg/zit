package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/alfa/toml"
	"github.com/friedenberg/zit/src/bravo/log"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/typ_akte"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/transacted"
	"github.com/friedenberg/zit/src/juliett/objekte"
	"github.com/friedenberg/zit/src/lima/objekte_store"
	"github.com/friedenberg/zit/src/mike/store_util"
)

type TypTransactedReader = objekte_store.TransactedReader[
	*kennung.Typ,
	sku.SkuLikePtr,
]

type typStore struct {
	*store_util.CommonStore[
		typ_akte.V0,
		*typ_akte.V0,
		kennung.Typ,
		*kennung.Typ,
	]
}

func makeTypStore(
	sa store_util.StoreUtil,
) (s *typStore, err error) {
	s = &typStore{}

	s.CommonStore, err = store_util.MakeCommonStore[
		typ_akte.V0,
		*typ_akte.V0,
		kennung.Typ,
		*kennung.Typ,
	](
		gattung.Typ,
		s,
		sa,
		s,
		objekte_store.MakeAkteFormat[typ_akte.V0, *typ_akte.V0](
			objekte.MakeTextParserIgnoreTomlErrors[typ_akte.V0](sa),
			objekte.ParsedAkteTomlFormatter[typ_akte.V0]{},
			sa,
		),
	)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	newOrUpdated := func(t *transacted.Typ) (err error) {
		s.StoreUtil.CommitUpdatedTransacted(t)

		if err = s.StoreUtil.GetKonfigPtr().AddTyp(t); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	s.CommonStore.CreateOrUpdater = objekte_store.MakeCreateOrUpdate[
		typ_akte.V0,
		*typ_akte.V0,
		kennung.Typ,
		*kennung.Typ,
	](
		sa,
		sa.GetLockSmith(),
		s.CommonStore,
		sa,
		TypTransactedReader(s),
		objekte_store.CreateOrUpdateDelegate[*transacted.Typ]{
			New: func(t *transacted.Typ) (err error) {
				if err = newOrUpdated(t); err != nil {
					err = errors.Wrap(err)
					return
				}

				return s.LogWriter.New(t)
			},
			Updated: func(t *transacted.Typ) (err error) {
				if err = newOrUpdated(t); err != nil {
					err = errors.Wrap(err)
					return
				}

				return s.LogWriter.Updated(t)
			},
			Unchanged: func(t *transacted.Typ) (err error) {
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

func (s typStore) Flush() (err error) {
	return
}

func (s typStore) AddOne(t *transacted.Typ) (err error) {
	s.StoreUtil.GetKonfigPtr().AddTyp(t)
	return
}

func (s typStore) UpdateOne(t *transacted.Typ) (err error) {
	log.Log().Printf("adding one: %s", t.GetSkuLike())
	s.StoreUtil.GetKonfigPtr().AddTyp(t)
	log.Log().Printf("done adding one: %s", t.GetSkuLike())
	return
}

// TODO-P3
func (s typStore) ReadAllSchwanzen(
	f schnittstellen.FuncIter[sku.SkuLikePtr],
) (err error) {
	// TODO-P2 switch to pointers
	if err = s.StoreUtil.GetKonfig().Typen.EachPtr(
		func(e *sku.Transacted2) (err error) {
			return f(e)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s typStore) ReadAll(
	f schnittstellen.FuncIter[sku.SkuLikePtr],
) (err error) {
	eachSku := func(sk sku.SkuLikePtr) (err error) {
		if sk.GetGattung() != gattung.Typ {
			return
		}

		var te *transacted.Typ

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

	if err = s.StoreUtil.GetBestandsaufnahmeStore().ReadAllSkus(eachSku); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s typStore) ReadOne(
	k1 schnittstellen.StringerGattungGetter,
) (tt *sku.Transacted2, err error) {
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