package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/toml"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/foxtrot/kennung"
	"github.com/friedenberg/zit/src/golf/age_io"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/golf/transaktion"
	"github.com/friedenberg/zit/src/hotel/objekte"
	"github.com/friedenberg/zit/src/india/typ"
)

type typLogWriter = collections.WriterFunc[*typ.Transacted]

type TypLogWriters struct {
	New, Updated, Unchanged typLogWriter
}

type typStore struct {
	common *common

	pool collections.PoolLike[typ.Transacted]

	objekte.TransactedInflator[
		typ.Objekte,
		*typ.Objekte,
		kennung.Typ,
		*kennung.Typ,
		objekte.NilVerzeichnisse[typ.Objekte],
		*objekte.NilVerzeichnisse[typ.Objekte],
	]

	objekte.AkteTextSaver[
		typ.Objekte,
		*typ.Objekte,
	]

	TypLogWriters
}

func (s *typStore) SetTypLogWriters(
	tlw TypLogWriters,
) {
	s.TypLogWriters = tlw
}

func makeTypStore(
	sa *common,
) (s *typStore, err error) {
	pool := collections.MakePool[typ.Transacted]()

	s = &typStore{
		common: sa,
		pool:   pool,
		TransactedInflator: objekte.MakeTransactedInflator[
			typ.Objekte,
			*typ.Objekte,
			kennung.Typ,
			*kennung.Typ,
			objekte.NilVerzeichnisse[typ.Objekte],
			*objekte.NilVerzeichnisse[typ.Objekte],
		](
			sa.ReadCloserObjektenSku,
			sa.AkteReader,
			nil,
			gattung.Parser[typ.Objekte, *typ.Objekte](
				typ.MakeFormatTextIgnoreTomlErrors(sa),
			),
			pool,
		),
		AkteTextSaver: objekte.MakeAkteTextSaver[
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

	if mutter != nil && tt.ObjekteSha().Equals(mutter.ObjekteSha()) {
		tt = mutter

		if err = s.TypLogWriters.Unchanged(tt); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	s.common.Transaktion.Skus.Add(&tt.Sku)
	s.common.KonfigPtr().AddTyp(tt)

	if mutter == nil {
		if err = s.TypLogWriters.New(tt); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if err = s.TypLogWriters.Updated(tt); err != nil {
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
	//TODO-P2 move to construction of inflator
	// p := collections.MakePool[*typ.Transacted]()

	if err = s.common.ReadAllTransaktions(
		func(t *transaktion.Transaktion) (err error) {
			if err = t.Skus.Each(
				func(o sku.SkuLike) (err error) {
					if o.GetGattung() != gattung.Typ {
						return
					}

					var te *typ.Transacted

					if te, err = s.Inflate(t.Time, o); err != nil {
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

					// if err = p.Apply(f, te); err != nil {
					// 	err = errors.Wrap(err)
					// 	return
					// }

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

func (s *typStore) reindexOne(
	t *transaktion.Transaktion,
	o sku.SkuLike,
) (err error) {
	var te *typ.Transacted
	defer s.pool.Put(te)

	if te, err = s.Inflate(t.Time, o); err != nil {
		if errors.Is(err, toml.Error{}) {
			err = nil
			return
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	s.common.KonfigPtr().AddTyp(te)

	if te.IsNew() {
		s.TypLogWriters.New(te)
	} else {
		s.TypLogWriters.Updated(te)
	}

	return
}
