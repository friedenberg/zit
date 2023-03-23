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
	"github.com/friedenberg/zit/src/hotel/objekte_store"
	"github.com/friedenberg/zit/src/hotel/typ"
	"github.com/friedenberg/zit/src/india/bestandsaufnahme"
	"github.com/friedenberg/zit/src/kilo/store_util"
)

type TypStore interface {
	CommonStore[
		typ.Objekte,
		*typ.Objekte,
		kennung.Typ,
		*kennung.Typ,
		objekte.NilVerzeichnisse[typ.Objekte],
		*objekte.NilVerzeichnisse[typ.Objekte],
	]
}

type TypTransactedReader = objekte_store.TransactedReader[
	*kennung.Typ,
	*typ.Transacted,
]

type typStore struct {
	*commonStore[
		typ.Objekte,
		*typ.Objekte,
		kennung.Typ,
		*kennung.Typ,
		objekte.NilVerzeichnisse[typ.Objekte],
		*objekte.NilVerzeichnisse[typ.Objekte],
	]

	objekte_store.CreateOrUpdater[
		*typ.Objekte,
		*kennung.Typ,
		*typ.Transacted,
		*typ.CheckedOut,
	]
}

func makeTypStore(
	sa store_util.StoreUtil,
) (s *typStore, err error) {
	s = &typStore{}

	s.commonStore, err = makeCommonStore[
		typ.Objekte,
		*typ.Objekte,
		kennung.Typ,
		*kennung.Typ,
		objekte.NilVerzeichnisse[typ.Objekte],
		*objekte.NilVerzeichnisse[typ.Objekte],
	](
		s,
		sa,
		s,
		nil,
		typ.MakeFormatTextIgnoreTomlErrors(sa),
		&typ.FormatterAkteTextToml{},
	)

	if err != nil {
		err = errors.Wrap(err)
		return
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
	s.StoreUtil.GetKonfigPtr().AddTyp(t)
	return
}

// TODO-P3
func (s typStore) ReadAllSchwanzen(
	f schnittstellen.FuncIter[*typ.Transacted],
) (err error) {
	if err = s.StoreUtil.GetKonfig().Typen.Each(f); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s typStore) ReadAll(
	f schnittstellen.FuncIter[*typ.Transacted],
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
	errors.TodoP3("inherited-typen-etiketten")
	tt = s.StoreUtil.GetKonfig().GetApproximatedTyp(*k).ActualOrNil()

	if tt == nil {
		err = errors.Wrap(objekte_store.ErrNotFound{Id: k})
		return
	}

	return
}

// func (s *typStore) Inherit(t *typ.Transacted) (err error) {
// 	if t == nil {
// 		panic("trying to inherit nil Typ")
// 	}

// 	errors.Log().Printf("inheriting %s", t.Sku.ObjekteSha)

// 	s.StoreUtil.CommitTransacted(t)
// 	old := s.StoreUtil.GetKonfig().GetApproximatedTyp(t.Sku.Kennung).ActualOrNil()

// 	if old == nil || old.Less(*t) {
// 		s.StoreUtil.GetKonfigPtr().AddTyp(t)
// 	}

// 	if t.IsNew() {
// 		s.LogWriter.New(t)
// 	} else {
// 		s.LogWriter.Updated(t)
// 	}

// 	return
// }
