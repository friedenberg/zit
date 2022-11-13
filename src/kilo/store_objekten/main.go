package store_objekten

import (
	"io"
	"os"
	"reflect"
	"sort"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/age"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/charlie/id"
	"github.com/friedenberg/zit/src/charlie/konfig"
	"github.com/friedenberg/zit/src/charlie/ts"
	"github.com/friedenberg/zit/src/delta/age_io"
	"github.com/friedenberg/zit/src/delta/hinweisen"
	"github.com/friedenberg/zit/src/delta/objekte"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/delta/zettel"
	"github.com/friedenberg/zit/src/echo/transaktion"
	"github.com/friedenberg/zit/src/foxtrot/zettel_named"
	"github.com/friedenberg/zit/src/hotel/zettel_transacted"
	"github.com/friedenberg/zit/src/india/zettel_verzeichnisse"
	"github.com/friedenberg/zit/src/juliett/store_verzeichnisse"
)

type LockSmith interface {
	IsAcquired() bool
}

type Store struct {
	lockSmith   LockSmith
	konfig      konfig.Konfig
	protoZettel zettel.ProtoZettel
	standort    standort.Standort
	age         age.Age

	zettelTransactedWriter collections.WriterFunc[*zettel_transacted.Zettel]

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

	if err = s.protoZettel.Typ.Set(k.Compiled.DefaultTyp); err != nil {
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

func (s *Store) SetZettelTransactedWriter(
	ztw collections.WriterFunc[*zettel_transacted.Zettel],
) {
	s.zettelTransactedWriter = ztw
}

func (s Store) Hinweisen() *hinweisen.Hinweisen {
	return s.hinweisen
}

func (s Store) WriteZettelObjekte(z zettel.Zettel) (sh sha.Sha, err error) {
	//no lock required

	var w *age_io.Mover

	mo := age_io.MoveOptions{
		Age:                      s.age,
		FinalPath:                s.standort.DirObjektenZettelen(),
		GenerateFinalPathFromSha: true,
	}

	if w, err = age_io.NewMover(mo); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer w.Close()

	c := zettel.FormatContextWrite{
		Zettel: z,
		Out:    w,
	}

	f := zettel.Objekte{}

	if _, err = f.WriteTo(c); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh = w.Sha()

	return
}

func (s Store) writeNamedZettelToIndex(tz zettel_transacted.Zettel) (err error) {
	if !s.lockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "write named zettel to index",
		}

		return
	}

	errors.Printf("writing zettel to index: %s", tz.Named)

	if err = s.verzeichnisseSchwanzen.Add(tz, tz.Named.Hinweis.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.verzeichnisseAll.Add(tz, tz.Named.Stored.Sha.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.indexKennung.addHinweis(tz.Named.Hinweis); err != nil {
		if errors.Is(err, hinweisen.ErrDoesNotExist{}) {
			errors.PrintErrf("kennung does not contain value: %s", err)
			err = nil
		} else {
			err = errors.Wrapf(err, "failed to write zettel to index: %s", tz.Named)
			return
		}
	}

	if err = s.indexAbbr.addZettelTransacted(tz); err != nil {
		err = errors.Wrapf(err, "failed to write zettel to index: %s", tz.Named)
		return
	}

	return
}

func (s Store) ReadHinweisSchwanzen(
	h hinweis.Hinweis,
) (zv zettel_transacted.Zettel, err error) {
	return s.verzeichnisseSchwanzen.ReadHinweisSchwanzen(h)
}

func (i *Store) ReadAllSchwanzenVerzeichnisse(
	ws ...collections.WriterFunc[*zettel_verzeichnisse.Zettel],
) (err error) {
	return i.verzeichnisseSchwanzen.ReadMany(ws...)
}

func (s Store) ReadAllSchwanzenTransacted(
	ws ...collections.WriterFunc[*zettel_transacted.Zettel],
) (err error) {
	w := zettel_verzeichnisse.MakeWriterZettelTransacted(
		collections.MakeChain(ws...),
	)

	return s.ReadAllSchwanzenVerzeichnisse(w)
}

func (i *Store) ReadAllVerzeichnisse(
	ws ...collections.WriterFunc[*zettel_verzeichnisse.Zettel],
) (err error) {
	return i.verzeichnisseAll.ReadMany(ws...)
}

func (s Store) ReadAllTransacted(
	ws ...collections.WriterFunc[*zettel_transacted.Zettel],
) (err error) {
	w := zettel_verzeichnisse.MakeWriterZettelTransacted(
		collections.MakeChain(ws...),
	)

	return s.ReadAllVerzeichnisse(w)
}

func (s Store) ReadOne(i id.Id) (tz zettel_transacted.Zettel, err error) {
	switch tid := i.(type) {
	case hinweis.Hinweis:
		if tz, err = s.verzeichnisseSchwanzen.ReadHinweisSchwanzen(tid); err != nil {
			err = errors.Wrap(err)
			return
		}

	default:
		err = errors.Errorf("unsupported identifier: %s, %#v", i, reflect.ValueOf(i))
	}

	return
}

func (s *Store) Create(in zettel.Zettel) (tz zettel_transacted.Zettel, err error) {
	if !s.lockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "create",
		}

		return
	}

	if in.IsEmpty() || s.protoZettel.Equals(in) {
		err = errors.Normalf("zettel is empty")
		return
	}

	s.protoZettel.Apply(&in)

	if err = in.ApplyKonfig(s.konfig); err != nil {
		err = errors.Wrap(err)
		return
	}

	tz.Named.Stored.Zettel = in

	if tz.Named.Stored.Sha, err = s.WriteZettelObjekte(tz.Named.Stored.Zettel); err != nil {
		err = errors.Wrap(err)
		return
	}

	//If the zettel exists, short circuit and return that
	// if tz2, err2 := s.Read(tz.Named.Stored.Sha); err2 == nil {
	// 	tz = tz2
	// 	return
	// }

	if tz.Named.Hinweis, err = s.indexKennung.createHinweis(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if tz, err = s.addZettelToTransaktion(tz.Named); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.writeNamedZettelToIndex(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.indexEtiketten.add(tz.Named.Stored.Zettel.Etiketten); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.zettelTransactedWriter(&tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) CreateWithHinweis(
	in zettel.Zettel,
	h hinweis.Hinweis,
) (tz zettel_transacted.Zettel, err error) {
	if !s.lockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "create with hinweis",
		}

		return
	}

	if in.IsEmpty() {
		err = errors.Normalf("zettel is empty")
		return
	}

	if err = in.ApplyKonfig(s.konfig); err != nil {
		err = errors.Wrap(err)
		return
	}

	tz.Named.Stored.Zettel = in

	if tz.Named.Stored.Sha, err = s.WriteZettelObjekte(tz.Named.Stored.Zettel); err != nil {
		err = errors.Wrap(err)
		return
	}

	if tz, err = s.addZettelToTransaktion(tz.Named); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.writeNamedZettelToIndex(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.indexEtiketten.add(tz.Named.Stored.Zettel.Etiketten); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.zettelTransactedWriter(&tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s Store) Etiketten() (es []etikett.Etikett, err error) {
	return s.indexEtiketten.allEtiketten()
}

// TODO support dry run
func (s *Store) Update(
	z *zettel_named.Zettel,
) (tz zettel_transacted.Zettel, err error) {
	if !s.lockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "update",
		}

		return
	}

	if err = z.Stored.Zettel.ApplyKonfig(s.konfig); err != nil {
		err = errors.Wrap(err)
		return
	}

	var mutter zettel_transacted.Zettel

	if mutter, err = s.verzeichnisseSchwanzen.ReadHinweisSchwanzen(
		z.Hinweis,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	tz.Named = *z

	if tz.Named.Stored.Sha, err = s.WriteZettelObjekte(z.Stored.Zettel); err != nil {
		err = errors.Wrap(err)
		return
	}

	if tz, err = s.addZettelToTransaktion(tz.Named); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.writeNamedZettelToIndex(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.indexEtiketten.addZettelWithOptionalMutter(&tz, &mutter); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.zettelTransactedWriter(&tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
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
		func(o *objekte.Objekte) (err error) {
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

			errors.Print(o)

			var chain []*zettel_transacted.Zettel

			if chain, err = s.AllInChain(*h); err != nil {
				err = errors.Wrap(err)
				return
			}

			var tz zettel_transacted.Zettel

			for _, someTz := range chain {
				errors.Print(someTz)
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

func (s Store) AllInChain(h hinweis.Hinweis) (c []*zettel_transacted.Zettel, err error) {
	mst := zettel_transacted.MakeMutableSetUnique(0)

	if err = s.verzeichnisseAll.ReadMany(
		func(z *zettel_verzeichnisse.Zettel) (err error) {
			if !z.Transacted.Named.Hinweis.Equals(h) {
				err = io.EOF
				return
			}

			return
		},
		zettel_verzeichnisse.MakeWriterZettelTransacted(mst.Add),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	c = mst.Elements()

	sort.Slice(
		c,
		func(i, j int) bool { return c[i].ObjekteTransacted().Less(c[j].ObjekteTransacted()) },
	)

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
		t.EachWithIndex(
			func(o *objekte.ObjekteWithIndex) (err error) {
				switch o.Gattung {

				case gattung.Zettel:
					var tz zettel_transacted.Zettel

					if tz, err = s.transactedZettelFromTransaktionObjekte(t, o); err != nil {
						if errors.Is(err, ErrNotFound{}) {
							errors.Print(err)
							err = nil
							return
						} else {
							err = errors.Wrap(err)
							return
						}
					}

					var mutter *zettel_transacted.Zettel

					if mutter1, err := s.verzeichnisseSchwanzen.ReadHinweisSchwanzen(tz.Named.Hinweis); err == nil {
						mutter = &mutter1
					}

					if err = s.writeNamedZettelToIndex(tz); err != nil {
						err = errors.Wrap(err)
						return
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
		)

		return
	}

	if err = s.ReadAllTransaktions(f); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
