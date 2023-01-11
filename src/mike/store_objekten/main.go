package store_objekten

import (
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/charlie/age"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/echo/standort"
	"github.com/friedenberg/zit/src/foxtrot/hinweis"
	"github.com/friedenberg/zit/src/foxtrot/ts"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/golf/transaktion"
	"github.com/friedenberg/zit/src/juliett/konfig_compiled"
	"github.com/friedenberg/zit/src/kilo/zettel"
)

type shaAbbr = sha.Abbr
type hinweisAbbr = hinweis.Abbr

type reindexer interface {
	reindexOne(*transaktion.Transaktion, sku.SkuLike) error
}

type Store struct {
	common

	//TODO move to methods
	shaAbbr
	hinweisAbbr

	zettelStore  *zettelStore
	typStore     TypStore
	etikettStore EtikettStore
	konfigStore  KonfigStore
}

func Make(
	lockSmith LockSmith,
	a age.Age,
	k *konfig_compiled.Compiled,
	st standort.Standort,
	p *collections.Pool[zettel.Transacted],
) (s *Store, err error) {
	s = &Store{
		common: common{
			LockSmith: lockSmith,
			Age:       a,
			konfig:    k,
			Standort:  st,
		},
	}

	t := ts.Now()

	for {
		p := s.TransaktionPath(t)

		if !files.Exists(p) {
			break
		}

		t.MoveForwardIota()
	}

	s.common.Transaktion = transaktion.MakeTransaktion(t)

	if s.common.Abbr, err = newIndexAbbr(
		&s.common,
		st.DirVerzeichnisse("Abbr"),
	); err != nil {
		err = errors.Wrapf(err, "failed to init abbr index")
		return
	}

	s.shaAbbr = s.common.Abbr
	s.hinweisAbbr = s.common.Abbr

	if s.zettelStore, err = makeZettelStore(&s.common, p); err != nil {
		err = errors.Wrap(err)
		return
	}

	if s.typStore, err = makeTypStore(&s.common); err != nil {
		err = errors.Wrap(err)
		return
	}

	if s.etikettStore, err = makeEtikettStore(&s.common); err != nil {
		err = errors.Wrap(err)
		return
	}

	if s.konfigStore, err = makeKonfigStore(&s.common); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) Abbr() *indexAbbr {
	return s.common.Abbr
}

func (s *Store) Zettel() *zettelStore {
	return s.zettelStore
}

func (s *Store) Typ() TypStore {
	return s.typStore
}

func (s *Store) Etikett() EtikettStore {
	return s.etikettStore
}

func (s *Store) Konfig() KonfigStore {
	return s.konfigStore
}

func (s *Store) CurrentTransaktionTime() ts.Time {
	return s.common.Transaktion.Time
}

func (s Store) RevertTransaktion(
	t *transaktion.Transaktion,
) (tzs zettel.MutableSet, err error) {
	if !s.common.LockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "revert",
		}

		return
	}

	tzs = zettel.MakeMutableSetUnique(t.Skus.Len())

	t.Skus.Each(
		func(o sku.SkuLike) (err error) {
			var h *hinweis.Hinweis
			ok := false

			if h, ok = o.GetId().(*hinweis.Hinweis); !ok {
				//TODO
				return
			}

			if !o.GetMutter()[1].IsZero() {
				err = errors.Errorf("merge reverts are not yet supported: %s", o)
				return
			}

			errors.Log().Print(o)

			var chain []*zettel.Transacted

			if chain, err = s.zettelStore.AllInChain(*h); err != nil {
				err = errors.Wrap(err)
				return
			}

			var tz *zettel.Transacted

			for _, someTz := range chain {
				errors.Log().Print(someTz)
				if someTz.Sku.Schwanz == o.GetMutter()[0] {
					tz = someTz
					break
				}
			}

			if tz.Sku.ObjekteSha.IsNull() {
				err = errors.Errorf("zettel not found in index!: %#v", o)
				return
			}

			if tz, err = s.zettelStore.Update(
				&tz.Objekte,
				&tz.Sku.Kennung,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			tzs.Add(tz)

			return
		},
	)

	return
}

func (s Store) Flush() (err error) {
	if !s.common.LockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "flush",
		}

		return
	}

	if s.common.Konfig().DryRun {
		return
	}

	if err = s.writeTransaktion(); err != nil {
		err = errors.Wrapf(err, "failed to write transaction")
		return
	}

	if err = s.zettelStore.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	// if err = s.typStore.Flush(); err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	// if err = s.etikettStore.Flush(); err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	if err = s.common.Abbr.Flush(); err != nil {
		errors.Err().Print(err)
		err = errors.Wrapf(err, "failed to flush abbr index")
		return
	}

	return
}

func (s *Store) Reindex() (err error) {
	if !s.common.LockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "reindex",
		}

		return
	}

	if err = os.RemoveAll(s.common.Standort.DirVerzeichnisse()); err != nil {
		err = errors.Wrapf(err, "failed to remove verzeichnisse dir")
		return
	}

	if err = os.MkdirAll(s.common.Standort.DirVerzeichnisse(), os.ModeDir|0755); err != nil {
		err = errors.Wrapf(err, "failed to make verzeichnisse dir")
		return
	}

	if err = os.MkdirAll(s.common.Standort.DirVerzeichnisseZettelenNeue(), os.ModeDir|0755); err != nil {
		err = errors.Wrapf(err, "failed to make verzeichnisse dir")
		return
	}

	if err = os.MkdirAll(s.common.Standort.DirVerzeichnisseZettelenNeueSchwanzen(), os.ModeDir|0755); err != nil {
		err = errors.Wrapf(err, "failed to make verzeichnisse dir")
		return
	}

	//TODO-P3 move to zettelStore
	if err = s.zettelStore.indexKennung.reset(); err != nil {
		err = errors.Wrapf(err, "failed to reset index kennung")
		return
	}

	f := func(t *transaktion.Transaktion) (err error) {
		errors.Out().Printf("%s/%s: %s", t.Time.Kopf(), t.Time.Schwanz(), t.Time)

		if err = t.Skus.Each(
			func(o sku.SkuLike) (err error) {
				switch o.GetGattung() {

				case gattung.Konfig:
					if err = s.konfigStore.reindexOne(
						t,
						o,
					); err != nil {
						err = errors.Wrap(err)
						return
					}

				case gattung.Typ:
					if err = s.typStore.reindexOne(
						t,
						o,
					); err != nil {
						err = errors.Wrapf(err, "Kennung: %s", o.GetId())
						return
					}

				case gattung.Etikett:
					if err = s.etikettStore.reindexOne(
						t,
						o,
					); err != nil {
						err = errors.Wrapf(
							err,
							"Sku: %s",
							o,
						)

						return
					}

				case gattung.Zettel:
					if err = s.zettelStore.reindexOne(
						t,
						o,
					); err != nil {
						err = errors.Wrap(err)
						return
					}

				default:
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
	}

	if err = s.ReadAllTransaktions(f); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
