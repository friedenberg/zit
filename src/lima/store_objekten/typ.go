package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/alfa/toml"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/golf/transaktion"
	"github.com/friedenberg/zit/src/hotel/bestandsaufnahme"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
	"github.com/friedenberg/zit/src/hotel/typ"
	"github.com/friedenberg/zit/src/store_util"
)

type TypStore interface {
	reindexer
	GattungStore

	objekte_store.Inheritor[*typ.Transacted]
	objekte_store.TransactedLogger[*typ.Transacted]

	objekte_store.AkteTextSaver[
		typ.Objekte,
		*typ.Objekte,
	]

	objekte_store.TransactedReader[
		*kennung.Typ,
		*typ.Transacted,
	]

	objekte_store.CreateOrUpdater[
		*typ.Objekte,
		*kennung.Typ,
		*typ.Transacted,
	]

	objekte_store.TransactedInflator[
		typ.Objekte,
		*typ.Objekte,
		kennung.Typ,
		*kennung.Typ,
		objekte.NilVerzeichnisse[typ.Objekte],
		*objekte.NilVerzeichnisse[typ.Objekte],
	]
}

type TypInflator = objekte_store.TransactedInflator[
	typ.Objekte,
	*typ.Objekte,
	kennung.Typ,
	*kennung.Typ,
	objekte.NilVerzeichnisse[typ.Objekte],
	*objekte.NilVerzeichnisse[typ.Objekte],
]

type TypTransactedReader = objekte_store.TransactedReader[
	*kennung.Typ,
	*typ.Transacted,
]

type TypLogWriter = objekte_store.LogWriter[*typ.Transacted]

type TypAkteTextSaver = objekte_store.AkteTextSaver[
	typ.Objekte,
	*typ.Objekte,
]

type typStore struct {
	store_util.StoreUtil

	pool collections.PoolLike[typ.Transacted]

	TypInflator
	TypAkteTextSaver
	TypLogWriter

	objekte_store.CreateOrUpdater[
		*typ.Objekte,
		*kennung.Typ,
		*typ.Transacted,
	]
}

func (s *typStore) SetLogWriter(
	tlw TypLogWriter,
) {
	s.TypLogWriter = tlw
}

func makeTypStore(
	sa store_util.StoreUtil,
) (s *typStore, err error) {
	pool := collections.MakePool[typ.Transacted]()

	s = &typStore{
		StoreUtil: sa,
		pool:                 pool,
		TypInflator: objekte_store.MakeTransactedInflator[
			typ.Objekte,
			*typ.Objekte,
			kennung.Typ,
			*kennung.Typ,
			objekte.NilVerzeichnisse[typ.Objekte],
			*objekte.NilVerzeichnisse[typ.Objekte],
		](
			sa,
			sa,
			nil,
			schnittstellen.Format[typ.Objekte, *typ.Objekte](
				typ.MakeFormatTextIgnoreTomlErrors(sa),
			),
			pool,
		),
		TypAkteTextSaver: objekte_store.MakeAkteTextSaver[
			typ.Objekte,
			*typ.Objekte,
		](
			sa,
			&typ.FormatterAkteTextToml{},
		),
	}

	newOrUpdated := func(t *typ.Transacted) (err error) {
		s.StoreUtil.CommitTransacted(t)
		s.StoreUtil.GetKonfigPtr().AddTyp(t)

		return
	}

	s.CreateOrUpdater = objekte_store.MakeCreateOrUpdate(
		sa,
		sa.GetLockSmith(),
		sa,
		TypTransactedReader(s),
		objekte_store.CreateOrUpdateDelegate[*typ.Transacted]{
			New: func(t *typ.Transacted) (err error) {
				if err = newOrUpdated(t); err != nil {
					err = errors.Wrap(err)
					return
				}

				return s.TypLogWriter.New(t)
			},
			Updated: func(t *typ.Transacted) (err error) {
				if err = newOrUpdated(t); err != nil {
					err = errors.Wrap(err)
					return
				}

				return s.TypLogWriter.Updated(t)
			},
			Unchanged: func(t *typ.Transacted) (err error) {
				return s.TypLogWriter.Unchanged(t)
			},
		},
	)

	return
}

func (s typStore) Flush() (err error) {
	return
}

// TODO-P3
func (s typStore) ReadAllSchwanzen(
	f collections.WriterFunc[*typ.Transacted],
) (err error) {
	if err = s.StoreUtil.GetKonfig().Typen.Each(f); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s typStore) ReadAll(
	f collections.WriterFunc[*typ.Transacted],
) (err error) {
	if s.StoreUtil.GetKonfig().UseBestandsaufnahme {
		f1 := func(t *bestandsaufnahme.Objekte) (err error) {
			if err = t.Akte.Skus.Each(
				func(sk sku.Sku2) (err error) {
					if sk.GetGattung() != gattung.Typ {
						return
					}

					var te *typ.Transacted

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
						if o.GetGattung() != gattung.Typ {
							return
						}

						var te *typ.Transacted

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

func (s typStore) ReadOne(
	k *kennung.Typ,
) (tt *typ.Transacted, err error) {
	errors.TodoP3("add support for working directory")
	at := s.StoreUtil.GetKonfig().GetApproximatedTyp(*k)

	if !at.HasValue() {
		err = errors.Wrap(objekte_store.ErrNotFound{Id: k})
		return
	}

	tt = at.Unwrap()

	return
}

func (s *typStore) Inherit(t *typ.Transacted) (err error) {
	if t == nil {
		panic("trying to inherit nil Typ")
	}

	errors.Log().Printf("inheriting %s", t.Sku.ObjekteSha)

	s.StoreUtil.CommitTransacted(t)
	old := s.StoreUtil.GetKonfig().GetApproximatedTyp(t.Sku.Kennung).ActualOrNil()

	if old == nil || old.Less(*t) {
		s.StoreUtil.GetKonfigPtr().AddTyp(t)
	}

	if t.IsNew() {
		s.TypLogWriter.New(t)
	} else {
		s.TypLogWriter.Updated(t)
	}

	return
}

func (s *typStore) reindexOne(
	sk sku.DataIdentity,
) (o schnittstellen.Stored, err error) {
	var te *typ.Transacted
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

	s.StoreUtil.GetKonfigPtr().AddTyp(te)

	if te.IsNew() {
		s.TypLogWriter.New(te)
	} else {
		s.TypLogWriter.Updated(te)
	}

	return
}
