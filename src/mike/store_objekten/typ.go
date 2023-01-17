package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/toml"
	"github.com/friedenberg/zit/src/bestandsaufnahme"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/foxtrot/kennung"
	"github.com/friedenberg/zit/src/golf/age_io"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/golf/transaktion"
	"github.com/friedenberg/zit/src/hotel/objekte"
	"github.com/friedenberg/zit/src/india/typ"
)

type TypStore interface {
	reindexer
	GattungStore

	objekte.Inheritor[*typ.Transacted]
	objekte.TransactedLogger[*typ.Transacted]

	objekte.AkteTextSaver[
		typ.Objekte,
		*typ.Objekte,
	]

	objekte.TransactedReader[
		*kennung.Typ,
		*typ.Transacted,
	]

	objekte.CreateOrUpdater[
		*typ.Objekte,
		*kennung.Typ,
		*typ.Transacted,
	]

	objekte.TransactedInflator[
		typ.Objekte,
		*typ.Objekte,
		kennung.Typ,
		*kennung.Typ,
		objekte.NilVerzeichnisse[typ.Objekte],
		*objekte.NilVerzeichnisse[typ.Objekte],
	]
}

type TypInflator = objekte.TransactedInflator[
	typ.Objekte,
	*typ.Objekte,
	kennung.Typ,
	*kennung.Typ,
	objekte.NilVerzeichnisse[typ.Objekte],
	*objekte.NilVerzeichnisse[typ.Objekte],
]

type TypLogWriter = objekte.LogWriter[*typ.Transacted]

type TypAkteTextSaver = objekte.AkteTextSaver[
	typ.Objekte,
	*typ.Objekte,
]

type typStore struct {
	common *common

	pool collections.PoolLike[typ.Transacted]

	TypInflator
	TypAkteTextSaver
	TypLogWriter
}

func (s *typStore) SetLogWriter(
	tlw TypLogWriter,
) {
	s.TypLogWriter = tlw
}

func makeTypStore(
	sa *common,
) (s *typStore, err error) {
	pool := collections.MakePool[typ.Transacted]()

	s = &typStore{
		common: sa,
		pool:   pool,
		TypInflator: objekte.MakeTransactedInflator[
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
			gattung.Format[typ.Objekte, *typ.Objekte](
				typ.MakeFormatTextIgnoreTomlErrors(sa),
			),
			pool,
		),
		TypAkteTextSaver: objekte.MakeAkteTextSaver[
			typ.Objekte,
			*typ.Objekte,
		](
			sa,
			&typ.FormatterAkteTextToml{},
		),
	}

	return
}

func (s typStore) Flush() (err error) {
	return
}

func (s typStore) CreateOrUpdate(
	to *typ.Objekte,
	tk *kennung.Typ,
) (tt *typ.Transacted, err error) {
	if !s.common.LockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "create or update typ",
		}

		return
	}

	var mutter *typ.Transacted

	if mutter, err = s.ReadOne(tk); err != nil {
		if errors.Is(err, ErrNotFound{}) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	tt = &typ.Transacted{
		Objekte: *to,
		Sku: sku.Transacted[kennung.Typ, *kennung.Typ]{
			Kennung: *tk,
			Verzeichnisse: sku.Verzeichnisse{
				Schwanz: s.common.Transaktion.Time,
			},
		},
	}

	//TODO-P3 refactor into reusable
	if mutter != nil {
		tt.Sku.Kopf = mutter.Sku.Kopf
		tt.Sku.Mutter[0] = mutter.Sku.Schwanz
	} else {
		tt.Sku.Kopf = s.common.Transaktion.Time
	}

	fo := objekte.MakeFormat[typ.Objekte, *typ.Objekte]()

	var w *age_io.Mover

	mo := age_io.MoveOptions{
		Age:                      s.common.Age,
		FinalPath:                s.common.Standort.DirObjektenTypen(),
		GenerateFinalPathFromSha: true,
	}

	if w, err = age_io.NewMover(mo); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, w.Close)

	if _, err = fo.Format(w, &tt.Objekte); err != nil {
		err = errors.Wrap(err)
		return
	}

	tt.Sku.ObjekteSha = w.Sha()

	if mutter != nil && tt.GetObjekteSha().Equals(mutter.GetObjekteSha()) {
		tt = mutter

		if err = s.TypLogWriter.Unchanged(tt); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	s.common.Transaktion.Skus.Add(&tt.Sku)
	s.common.KonfigPtr().AddTyp(tt)

	if mutter == nil {
		if err = s.TypLogWriter.New(tt); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if err = s.TypLogWriter.Updated(tt); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

// TODO-P3
func (s typStore) ReadAllSchwanzen(
	f collections.WriterFunc[*typ.Transacted],
) (err error) {
	if err = s.common.konfig.Typen.Each(f); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s typStore) ReadAll(
	f collections.WriterFunc[*typ.Transacted],
) (err error) {
	if s.common.Konfig().UseBestandsaufnahme {
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

		if err = s.common.bestandsaufnahmeStore.ReadAll(f1); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if err = s.common.ReadAllTransaktions(
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
	tt = s.common.Konfig().GetTyp(*k)

	if tt == nil {
		err = errors.Wrap(ErrNotFound{Id: k})
		return
	}

	return
}

func (s typStore) AllInChain(k kennung.Typ) (c []*typ.Transacted, err error) {
	return
}

func (s *typStore) Inherit(t *typ.Transacted) (err error) {
	errors.Log().Printf("inheriting %s", t.Sku.ObjekteSha)

	s.common.Bestandsaufnahme.Akte.Skus.Push(t.Sku.Sku2())
	s.common.Transaktion.Skus.Add(&t.Sku)

	if old := s.common.Konfig().GetTyp(t.Sku.Kennung); old == nil || old.Less(*t) {
		s.common.KonfigPtr().AddTyp(t)
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
) (o gattung.Stored, err error) {
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

	s.common.KonfigPtr().AddTyp(te)

	if te.IsNew() {
		s.TypLogWriter.New(te)
	} else {
		s.TypLogWriter.Updated(te)
	}

	return
}
