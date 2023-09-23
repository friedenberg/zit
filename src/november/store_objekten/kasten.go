package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/alfa/toml"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/kasten_akte"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/juliett/objekte"
	"github.com/friedenberg/zit/src/lima/objekte_store"
	"github.com/friedenberg/zit/src/mike/store_util"
)

type KastenTransactedReader = objekte_store.TransactedReader[
	*kennung.Kasten,
	sku.SkuLikePtr,
]

type kastenStore struct {
	*store_util.CommonStore[
		kasten_akte.V0,
		*kasten_akte.V0,
		kennung.Kasten,
		*kennung.Kasten,
	]
}

func makeKastenStore(
	sa store_util.StoreUtil,
) (s *kastenStore, err error) {
	s = &kastenStore{}

	s.CommonStore, err = store_util.MakeCommonStore[
		kasten_akte.V0,
		*kasten_akte.V0,
		kennung.Kasten,
		*kennung.Kasten,
	](
		gattung.Kasten,
		s,
		sa,
		s,
		objekte_store.MakeAkteFormat[kasten_akte.V0, *kasten_akte.V0](
			objekte.MakeTextParserIgnoreTomlErrors[kasten_akte.V0](sa),
			objekte.ParsedAkteTomlFormatter[kasten_akte.V0]{},
			sa,
		),
	)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	newOrUpdated := func(t *sku.Transacted) (err error) {
		s.StoreUtil.CommitUpdatedTransacted(t)

		if err = s.StoreUtil.GetKonfigPtr().AddKasten(t); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	s.CommonStore.CreateOrUpdater = objekte_store.MakeCreateOrUpdate[
		kasten_akte.V0,
		*kasten_akte.V0,
		kennung.Kasten,
		*kennung.Kasten,
	](
		sa,
		sa.GetLockSmith(),
		sa.ObjekteReaderWriterFactory(gattung.Kasten),
		sa,
		KastenTransactedReader(s),
		objekte_store.CreateOrUpdateDelegate[*sku.Transacted]{
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

func (s kastenStore) Flush() (err error) {
	return
}

func (s kastenStore) AddOne(t *sku.Transacted) (err error) {
	s.StoreUtil.GetKonfigPtr().AddKasten(t)
	return
}

func (s kastenStore) UpdateOne(t *sku.Transacted) (err error) {
	s.StoreUtil.GetKonfigPtr().AddKasten(t)
	return
}

// TODO-P3
func (s kastenStore) ReadAllSchwanzen(
	f schnittstellen.FuncIter[sku.SkuLikePtr],
) (err error) {
	if err = s.StoreUtil.GetKonfig().Kisten.EachPtr(
		func(e *sku.Transacted) (err error) {
			return f(e)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s kastenStore) ReadAll(
	f schnittstellen.FuncIter[sku.SkuLikePtr],
) (err error) {
	eachSku := func(sk sku.SkuLikePtr) (err error) {
		if sk.GetGattung() != gattung.Kasten {
			return
		}

		var te *sku.Transacted

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
