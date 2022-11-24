package store_objekten

import (
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/age"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/delta/hinweisen"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/delta/ts"
	"github.com/friedenberg/zit/src/echo/konfig"
	"github.com/friedenberg/zit/src/echo/sku"
	"github.com/friedenberg/zit/src/golf/transaktion"
	"github.com/friedenberg/zit/src/india/zettel"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
	"github.com/friedenberg/zit/src/juliett/zettel_verzeichnisse"
	"github.com/friedenberg/zit/src/kilo/store_verzeichnisse"
)

type Store struct {
	lockSmith   LockSmith
	konfig      konfig.Konfig
	protoZettel zettel.ProtoZettel
	standort    standort.Standort
	age         age.Age

	zettelTransactedWriter ZettelTransactedLogWriters

	hinweisen *hinweisen.Hinweisen
	// *indexZettelen
	*indexEtiketten
	*indexKennung
	*indexAbbr

	verzeichnisseSchwanzen *verzeichnisseSchwanzen
	verzeichnisseAll       *store_verzeichnisse.Zettelen

	transaktion.Transaktion
}

func Make(
	lockSmith LockSmith,
	a age.Age,
	k konfig.Konfig,
	st standort.Standort,
	p zettel_verzeichnisse.Pool,
) (s *Store, err error) {
	s = &Store{
		lockSmith: lockSmith,
		age:       a,
		konfig:    k,
		standort:  st,
	}

	s.protoZettel = zettel.MakeProtoZettel()

	if err = s.protoZettel.Typ.Set(k.Compiled.DefaultTyp.Name.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if s.hinweisen, err = hinweisen.New(st.DirZit()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if s.verzeichnisseSchwanzen, err = makeVerzeichnisseSchwanzen(
		k,
		st,
		s,
		p,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if s.verzeichnisseAll, err = store_verzeichnisse.MakeZettelen(
		k,
		st.DirVerzeichnisseZettelenNeue(),
		s,
		p,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	// s.indexZettelen, err = newIndexZettelen(
	// 	st.FileVerzeichnisseZettelen(),
	// 	s,
	// )

	if err != nil {
		err = errors.Wrapf(err, "failed to init zettel index")
		return
	}

	s.indexEtiketten, err = newIndexEtiketten(
		st.FileVerzeichnisseEtiketten(),
		s,
	)

	if err != nil {
		err = errors.Wrapf(err, "failed to init zettel index")
		return
	}

	s.indexKennung, err = newIndexKennung(
		k,
		s,
		s.hinweisen,
		st.DirVerzeichnisse("Kennung"),
	)

	if err != nil {
		err = errors.Wrapf(err, "failed to init kennung index")
		return
	}

	s.indexAbbr, err = newIndexAbbr(
		s,
		st.DirVerzeichnisse("Abbr"),
	)

	if err != nil {
		err = errors.Wrapf(err, "failed to init abbr index")
		return
	}

	s.Transaktion = transaktion.MakeTransaktion(ts.Now())

	return
}

func (s *Store) CurrentTransaktionTime() ts.Time {
	return s.Transaktion.Time
}

func (s Store) Hinweisen() *hinweisen.Hinweisen {
	return s.hinweisen
}

func (s Store) RevertTransaktion(
	t *transaktion.Transaktion,
) (tzs zettel_transacted.MutableSet, err error) {
	if !s.lockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "revert",
		}

		return
	}

	tzs = zettel_transacted.MakeMutableSetUnique(t.Len())

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

			var chain []*zettel_transacted.Zettel

			if chain, err = s.AllInChain(*h); err != nil {
				err = errors.Wrap(err)
				return
			}

			var tz zettel_transacted.Zettel

			for _, someTz := range chain {
				errors.Log().Print(someTz)
				if someTz.Schwanz == o.Mutter[0] {
					tz = *someTz
					break
				}
			}

			if tz.Named.Stored.Sha.IsNull() {
				err = errors.Errorf("zettel not found in index!: %#v", o)
				return
			}

			if tz, err = s.Update(&tz.Named); err != nil {
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
	if !s.lockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "flush",
		}

		return
	}

	if s.konfig.DryRun {
		return
	}

	if err = s.writeTransaktion(); err != nil {
		err = errors.Wrapf(err, "failed to write transaction")
		return
	}

	if err = s.verzeichnisseSchwanzen.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.verzeichnisseAll.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	// if err = s.indexZettelen.Flush(); err != nil {
	// 	err = errors.Wrapf(err, "failed to flush new zettel index")
	// 	return
	// }

	if err = s.indexEtiketten.Flush(); err != nil {
		err = errors.Wrapf(err, "failed to flush new zettel index")
		return
	}

	if err = s.indexKennung.Flush(); err != nil {
		err = errors.Wrapf(err, "failed to flush new kennung index")
		return
	}

	if err = s.indexAbbr.Flush(); err != nil {
		err = errors.Wrapf(err, "failed to flush abbr index")
		return
	}

	return
}

func (s *Store) Reindex() (err error) {
	if !s.lockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "reindex",
		}

		return
	}

	if err = os.RemoveAll(s.standort.DirVerzeichnisse()); err != nil {
		err = errors.Wrapf(err, "failed to remove verzeichnisse dir")
		return
	}

	if err = os.MkdirAll(s.standort.DirVerzeichnisse(), os.ModeDir|0755); err != nil {
		err = errors.Wrapf(err, "failed to make verzeichnisse dir")
		return
	}

	if err = os.MkdirAll(s.standort.DirVerzeichnisseZettelenNeue(), os.ModeDir|0755); err != nil {
		err = errors.Wrapf(err, "failed to make verzeichnisse dir")
		return
	}

	if err = os.MkdirAll(s.standort.DirVerzeichnisseZettelenNeueSchwanzen(), os.ModeDir|0755); err != nil {
		err = errors.Wrapf(err, "failed to make verzeichnisse dir")
		return
	}

	if err = s.indexKennung.reset(); err != nil {
		err = errors.Wrapf(err, "failed to reset index kennung")
		return
	}

	f := func(t *transaktion.Transaktion) (err error) {
		if err = t.EachWithIndex(
			func(o *sku.Indexed) (err error) {
				switch o.Gattung {

				case gattung.Zettel:
					var tz zettel_transacted.Zettel

					if tz, err = s.transactedZettelFromTransaktionObjekte(t, o); err != nil {
						err = errors.Wrap(err)
						return
					}

					var mutter *zettel_transacted.Zettel

					if mutter1, err := s.verzeichnisseSchwanzen.ReadHinweisSchwanzen(tz.Named.Kennung); err == nil {
						mutter = &mutter1
					}

					if err = s.writeNamedZettelToIndex(tz); err != nil {
						err = errors.Wrap(err)
						return
					}

					if mutter == nil {
						if err = s.zettelTransactedWriter.New(&tz); err != nil {
							err = errors.Wrap(err)
							return
						}
					} else {
						if err = s.zettelTransactedWriter.Updated(&tz); err != nil {
							err = errors.Wrap(err)
							return
						}
					}

					if err = s.indexEtiketten.addZettelWithOptionalMutter(&tz, mutter); err != nil {
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
