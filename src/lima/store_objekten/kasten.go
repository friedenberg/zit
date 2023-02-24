package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/alfa/toml"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/transaktion"
	"github.com/friedenberg/zit/src/hotel/kasten"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
	"github.com/friedenberg/zit/src/india/bestandsaufnahme"
	"github.com/friedenberg/zit/src/kilo/store_util"
)

type KastenStore interface {
	reindexer
	GattungStore

	objekte_store.Inheritor[*kasten.Transacted]
	objekte_store.TransactedLogger[*kasten.Transacted]

	objekte_store.AkteTextSaver[
		kasten.Objekte,
		*kasten.Objekte,
	]

	objekte_store.TransactedReader[
		*kennung.Kasten,
		*kasten.Transacted,
	]

	objekte_store.CreateOrUpdater[
		*kasten.Objekte,
		*kennung.Kasten,
		*kasten.Transacted,
		*kasten.CheckedOut,
	]

	objekte_store.TransactedInflator[
		kasten.Objekte,
		*kasten.Objekte,
		kennung.Kasten,
		*kennung.Kasten,
		kasten.Verzeichnisse,
		*kasten.Verzeichnisse,
	]
}

type KastenInflator = objekte_store.TransactedInflator[
	kasten.Objekte,
	*kasten.Objekte,
	kennung.Kasten,
	*kennung.Kasten,
	kasten.Verzeichnisse,
	*kasten.Verzeichnisse,
]

type KastenTransactedReader = objekte_store.TransactedReader[
	*kennung.Kasten,
	*kasten.Transacted,
]

type KastenLogWriter = objekte_store.LogWriter[*kasten.Transacted]

type KastenAkteTextSaver = objekte_store.AkteTextSaver[
	kasten.Objekte,
	*kasten.Objekte,
]

type kastenStore struct {
	store_util.StoreUtil

	pool schnittstellen.Pool[kasten.Transacted, *kasten.Transacted]

	KastenInflator
	KastenAkteTextSaver
	KastenLogWriter

	objekte_store.CreateOrUpdater[
		*kasten.Objekte,
		*kennung.Kasten,
		*kasten.Transacted,
		*kasten.CheckedOut,
	]
}

func (s *kastenStore) SetLogWriter(
	tlw KastenLogWriter,
) {
	s.KastenLogWriter = tlw
}

func makeKastenStore(
	sa store_util.StoreUtil,
) (s *kastenStore, err error) {
	pool := collections.MakePool[kasten.Transacted]()

	s = &kastenStore{
		StoreUtil: sa,
		pool:      pool,
		KastenInflator: objekte_store.MakeTransactedInflator[
			kasten.Objekte,
			*kasten.Objekte,
			kennung.Kasten,
			*kennung.Kasten,
			kasten.Verzeichnisse,
			*kasten.Verzeichnisse,
		](
			sa,
			sa,
			nil,
			schnittstellen.Format[kasten.Objekte, *kasten.Objekte](
				kasten.MakeFormatTextIgnoreTomlErrors(sa),
			),
			pool,
		),
		KastenAkteTextSaver: objekte_store.MakeAkteTextSaver[
			kasten.Objekte,
			*kasten.Objekte,
		](
			sa,
			&kasten.FormatterAkteTextToml{},
		),
	}

	newOrUpdated := func(t *kasten.Transacted) (err error) {
		s.StoreUtil.CommitTransacted(t)
		s.StoreUtil.GetKonfigPtr().AddKasten(t)

		return
	}

	s.CreateOrUpdater = objekte_store.MakeCreateOrUpdate(
		sa,
		sa.GetLockSmith(),
		sa,
		KastenTransactedReader(s),
		objekte_store.CreateOrUpdateDelegate[*kasten.Transacted]{
			New: func(t *kasten.Transacted) (err error) {
				if err = newOrUpdated(t); err != nil {
					err = errors.Wrap(err)
					return
				}

				return s.KastenLogWriter.New(t)
			},
			Updated: func(t *kasten.Transacted) (err error) {
				if err = newOrUpdated(t); err != nil {
					err = errors.Wrap(err)
					return
				}

				return s.KastenLogWriter.Updated(t)
			},
			Unchanged: func(t *kasten.Transacted) (err error) {
				return s.KastenLogWriter.Unchanged(t)
			},
		},
	)

	return
}

func (s kastenStore) Flush() (err error) {
	return
}

// TODO-P3
func (s kastenStore) ReadAllSchwanzen(
	f schnittstellen.FuncIter[*kasten.Transacted],
) (err error) {
	if err = s.StoreUtil.GetKonfig().Kisten.Each(f); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s kastenStore) ReadAll(
	f schnittstellen.FuncIter[*kasten.Transacted],
) (err error) {
	if s.StoreUtil.GetKonfig().UseBestandsaufnahme {
		f1 := func(t *bestandsaufnahme.Objekte) (err error) {
			if err = t.Akte.Skus.Each(
				func(sk sku.Sku2) (err error) {
					if sk.GetGattung() != gattung.Kasten {
						return
					}

					var te *kasten.Transacted

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
					t.Tai,
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
						if o.GetGattung() != gattung.Kasten {
							return
						}

						var te *kasten.Transacted

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

func (s *kastenStore) Inherit(t *kasten.Transacted) (err error) {
	if t == nil {
		panic("trying to inherit nil Kasten")
	}

	errors.Log().Printf("inheriting %s", t.Sku.ObjekteSha)

	s.StoreUtil.CommitTransacted(t)
	old := s.StoreUtil.GetKonfig().GetKasten(t.Sku.Kennung)

	if old == nil || old.Less(*t) {
		s.StoreUtil.GetKonfigPtr().AddKasten(t)
	}

	if t.IsNew() {
		s.KastenLogWriter.New(t)
	} else {
		s.KastenLogWriter.Updated(t)
	}

	return
}

func (s *kastenStore) reindexOne(
	sk sku.DataIdentity,
) (o schnittstellen.Stored, err error) {
	var te *kasten.Transacted
	defer s.pool.Put(te)

	if te, err = s.InflateFromDataIdentity(sk); err != nil {
		if errors.Is(err, toml.Error{}) {
			err = nil
			return
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	o = te

	s.StoreUtil.GetKonfigPtr().AddKasten(te)

	if te.IsNew() {
		s.KastenLogWriter.New(te)
	} else {
		s.KastenLogWriter.Updated(te)
	}

	return
}
