package store_objekten

import (
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/charlie/age"
	"github.com/friedenberg/zit/src/charlie/gattung"
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

type Store struct {
	common common

	//TODO move to methods
	ioFactory
	shaAbbr
	hinweisAbbr

	zettelStore  *zettelStore
	typStore     *typStore
	etikettStore *etikettStore
	konfigStore  *konfigStore
}

func Make(
	lockSmith LockSmith,
	a age.Age,
	k konfig_compiled.Compiled,
	st standort.Standort,
	p *zettel.PoolVerzeichnisse,
) (s *Store, err error) {
	s = &Store{
		common: common{
			LockSmith: lockSmith,
			Age:       a,
			Konfig:    k,
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

	s.ioFactory = s.common

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

func (s *Store) Typ() *typStore {
	return s.typStore
}

func (s *Store) Etikett() *etikettStore {
	return s.etikettStore
}

func (s *Store) Konfig() *konfigStore {
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

	tzs = zettel.MakeMutableSetUnique(t.Len())

	t.Each(
		func(o *sku.Sku) (err error) {
			var h *hinweis.Hinweis
			ok := false

			if h, ok = o.Id.(*hinweis.Hinweis); !ok {
				//TODO
				return
			}

			if !o.Mutter[1].IsZero() {
				err = errors.Errorf("merge reverts are not yet supported: %s", o)
				return
			}

			errors.Log().Print(o)

			var chain []*zettel.Transacted

			if chain, err = s.zettelStore.AllInChain(*h); err != nil {
				err = errors.Wrap(err)
				return
			}

			var tz zettel.Transacted

			for _, someTz := range chain {
				errors.Log().Print(someTz)
				if someTz.Sku.Schwanz == o.Mutter[0] {
					tz = *someTz
					break
				}
			}

			if tz.Sku.Sha.IsNull() {
				err = errors.Errorf("zettel not found in index!: %#v", o)
				return
			}

			if tz, err = s.zettelStore.Update(&tz.Objekte, &tz.Sku.Kennung); err != nil {
				err = errors.Wrap(err)
				return
			}

			tzs.Add(&tz)

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

	if s.common.Konfig.DryRun {
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

	//TODO move all below to zettelStore
	if err = s.zettelStore.indexKennung.reset(); err != nil {
		err = errors.Wrapf(err, "failed to reset index kennung")
		return
	}

	f := func(t *transaktion.Transaktion) (err error) {
		if err = t.Each(
			func(o *sku.Sku) (err error) {
				switch o.Gattung {

				case gattung.Zettel:
					var tz zettel.Transacted

					if tz, err = s.zettelStore.transactedZettelFromTransaktionObjekte(t, o); err != nil {
						//TODO decide on how to handle format errors
						errors.Err().Print(err)
						err = nil
						// err = errors.Wrap(err)
						return
					}

					var mutter *zettel.Transacted

					if mutter1, err := s.zettelStore.verzeichnisseSchwanzen.ReadHinweisSchwanzen(tz.Sku.Kennung); err == nil {
						mutter = &mutter1
					}

					if err = s.zettelStore.writeNamedZettelToIndex(tz); err != nil {
						err = errors.Wrap(err)
						return
					}

					if mutter == nil {
						if err = s.zettelStore.zettelTransactedWriter.New(&tz); err != nil {
							err = errors.Wrap(err)
							return
						}
					} else {
						if err = s.zettelStore.zettelTransactedWriter.Updated(&tz); err != nil {
							err = errors.Wrap(err)
							return
						}
					}

					if err = s.zettelStore.indexEtiketten.addZettelWithOptionalMutter(&tz, mutter); err != nil {
						err = errors.Wrap(err)
						return
					}

				default:
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

	if err = s.ReadAllTransaktions(f); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}