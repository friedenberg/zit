package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/alfa/toml"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/etikett_akte"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/transacted"
	"github.com/friedenberg/zit/src/juliett/objekte"
	"github.com/friedenberg/zit/src/lima/objekte_store"
	"github.com/friedenberg/zit/src/mike/store_util"
)

type EtikettStore interface {
	CommonStore[
		etikett_akte.V0,
		*etikett_akte.V0,
		kennung.Etikett,
		*kennung.Etikett,
	]
}

type EtikettTransactedReader = objekte_store.TransactedReader[
	*kennung.Etikett,
	*transacted.Etikett,
]

type etikettStore struct {
	*commonStore[
		etikett_akte.V0,
		*etikett_akte.V0,
		kennung.Etikett,
		*kennung.Etikett,
	]
}

func makeEtikettStore(
	sa store_util.StoreUtil,
) (s *etikettStore, err error) {
	s = &etikettStore{}

	s.commonStore, err = makeCommonStore[
		etikett_akte.V0,
		*etikett_akte.V0,
		kennung.Etikett,
		*kennung.Etikett,
	](
		gattung.Etikett,
		s,
		sa,
		s,
		objekte_store.MakeAkteFormat[etikett_akte.V0, *etikett_akte.V0](
			objekte.MakeTextParserIgnoreTomlErrors[etikett_akte.V0](sa),
			objekte.ParsedAkteTomlFormatter[etikett_akte.V0]{},
			sa,
		),
	)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	newOrUpdated := func(t *transacted.Etikett) (err error) {
		s.StoreUtil.CommitUpdatedTransacted(t)

		if err = s.StoreUtil.GetKonfigPtr().AddEtikett(t); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	s.commonStore.CreateOrUpdater = objekte_store.MakeCreateOrUpdate[
		etikett_akte.V0,
		*etikett_akte.V0,
		kennung.Etikett,
		*kennung.Etikett,
	](
		sa,
		sa.GetLockSmith(),
		sa.ObjekteReaderWriterFactory(gattung.Etikett),
		sa,
		EtikettTransactedReader(s),
		objekte_store.CreateOrUpdateDelegate[*transacted.Etikett]{
			New: func(t *transacted.Etikett) (err error) {
				if err = newOrUpdated(t); err != nil {
					err = errors.Wrap(err)
					return
				}

				return s.LogWriter.New(t)
			},
			Updated: func(t *transacted.Etikett) (err error) {
				if err = newOrUpdated(t); err != nil {
					err = errors.Wrap(err)
					return
				}

				return s.LogWriter.Updated(t)
			},
			Unchanged: func(t *transacted.Etikett) (err error) {
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

func (s etikettStore) addOne(t *transacted.Etikett) (err error) {
	s.StoreUtil.GetKonfigPtr().AddEtikett(t)
	return
}

func (s etikettStore) updateOne(t *transacted.Etikett) (err error) {
	s.StoreUtil.GetKonfigPtr().AddEtikett(t)
	return
}

func (s etikettStore) ReadOne(
	k *kennung.Etikett,
) (tt *transacted.Etikett, err error) {
	tt1 := s.StoreUtil.GetKonfig().GetEtikett(*k)
	tt = &tt1

	if tt == nil {
		err = errors.Wrap(objekte_store.ErrNotFound{Id: k})
		return
	}

	return
}

func (s etikettStore) ReadAllSchwanzen(
	f schnittstellen.FuncIter[*transacted.Etikett],
) (err error) {
	if err = s.StoreUtil.GetKonfig().EachEtikett(f); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s etikettStore) ReadAll(
	f schnittstellen.FuncIter[*transacted.Etikett],
) (err error) {
	eachSku := func(o sku.SkuLikePtr) (err error) {
		if o.GetGattung() != gattung.Etikett {
			return
		}

		var te *transacted.Etikett

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

	if err = s.StoreUtil.GetBestandsaufnahmeStore().ReadAllSkus(eachSku); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
