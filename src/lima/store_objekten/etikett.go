package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/alfa/toml"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/golf/transaktion"
	"github.com/friedenberg/zit/src/hotel/etikett"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
	"github.com/friedenberg/zit/src/india/bestandsaufnahme"
	"github.com/friedenberg/zit/src/kilo/store_util"
)

type EtikettStore interface {
	CommonStore[
		etikett.Akte,
		*etikett.Akte,
		kennung.Etikett,
		*kennung.Etikett,
		objekte.NilVerzeichnisse[etikett.Akte],
		*objekte.NilVerzeichnisse[etikett.Akte],
	]
}

type EtikettTransactedReader = objekte_store.TransactedReader[
	*kennung.Etikett,
	*etikett.Transacted,
]

type etikettStore struct {
	*commonStore[
		etikett.Akte,
		*etikett.Akte,
		kennung.Etikett,
		*kennung.Etikett,
		objekte.NilVerzeichnisse[etikett.Akte],
		*objekte.NilVerzeichnisse[etikett.Akte],
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
		objekte.NilVerzeichnisse[etikett.Akte],
		*objekte.NilVerzeichnisse[etikett.Akte],
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

	newOrUpdated := func(t *etikett.Transacted) (err error) {
		s.StoreUtil.CommitUpdatedTransacted(t)
		s.StoreUtil.GetKonfigPtr().AddEtikett(t)

		return
	}

	s.commonStore.CreateOrUpdater = objekte_store.MakeCreateOrUpdate(
		sa,
		sa.GetLockSmith(),
		sa.ObjekteReaderWriterFactory(gattung.Etikett),
		sa,
		EtikettTransactedReader(s),
		objekte_store.CreateOrUpdateDelegate[*etikett.Transacted]{
			New: func(t *etikett.Transacted) (err error) {
				if err = newOrUpdated(t); err != nil {
					err = errors.Wrap(err)
					return
				}

				return s.LogWriter.New(t)
			},
			Updated: func(t *etikett.Transacted) (err error) {
				if err = newOrUpdated(t); err != nil {
					err = errors.Wrap(err)
					return
				}

				return s.LogWriter.Updated(t)
			},
			Unchanged: func(t *etikett.Transacted) (err error) {
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

func (s etikettStore) addOne(t *etikett.Transacted) (err error) {
	s.StoreUtil.GetKonfigPtr().AddEtikett(t)
	return
}

func (s etikettStore) updateOne(t *etikett.Transacted) (err error) {
	s.StoreUtil.GetKonfigPtr().AddEtikett(t)
	return
}

func (s etikettStore) ReadOne(
	k *kennung.Etikett,
) (tt *etikett.Transacted, err error) {
	tt = s.StoreUtil.GetKonfig().GetEtikett(*k)

	if tt == nil {
		err = errors.Wrap(objekte_store.ErrNotFound{Id: k})
		return
	}

	return
}

func (s etikettStore) ReadAllSchwanzen(
	f schnittstellen.FuncIter[*etikett.Transacted],
) (err error) {
	if err = s.StoreUtil.GetKonfig().Etiketten.Each(f); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s etikettStore) ReadAll(
	f schnittstellen.FuncIter[*etikett.Transacted],
) (err error) {
	if s.StoreUtil.GetKonfig().UseBestandsaufnahme {
		f1 := func(t *bestandsaufnahme.Transacted) (err error) {
			if err = t.Akte.Skus.Each(
				func(sk sku.Sku) (err error) {
					if sk.GetGattung() != gattung.Etikett {
						return
					}

					var te *etikett.Transacted

					if te, err = s.InflateFromDataIdentity(sk); err != nil {
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
				},
			); err != nil {
				err = errors.Wrapf(
					err,
					"Bestandsaufnahme: %s",
					t.GetKennung(),
				)

				return
			}

			return
		}

		if err = s.StoreUtil.GetBestandsaufnahmeStore().ReadAll(f1); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if err = s.StoreUtil.GetTransaktionStore().ReadAllTransaktions(
			func(t *transaktion.Transaktion) (err error) {
				if err = t.Skus.Each(
					func(o sku.SkuLike) (err error) {
						if o.GetGattung() != gattung.Etikett {
							return
						}

						var te *etikett.Transacted

						if te, err = s.InflateFromDataIdentity(o); err != nil {
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
					},
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
	}

	return
}
