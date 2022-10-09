package store_objekten

import (
	"os"
	"reflect"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/paper"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/age"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/charlie/id"
	"github.com/friedenberg/zit/src/charlie/konfig"
	"github.com/friedenberg/zit/src/charlie/ts"
	"github.com/friedenberg/zit/src/delta/age_io"
	"github.com/friedenberg/zit/src/delta/hinweisen"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/delta/transaktion"
	"github.com/friedenberg/zit/src/delta/zettel"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
	"github.com/friedenberg/zit/src/hotel/zettel_verzeichnisse"
	"github.com/friedenberg/zit/src/india/store_verzeichnisse"
)

type LockSmith interface {
	IsAcquired() bool
}

type ZettelTransactedPrinter interface {
	ZettelTransacted(zettel_transacted.Zettel) *paper.Paper
}

type Store struct {
	lockSmith LockSmith
	konfig    konfig.Konfig
	standort  standort.Standort
	age       age.Age

	zettelTransactedPrinter ZettelTransactedPrinter
	hinweisen               *hinweisen.Hinweisen
	*indexZettelen
	// *indexZettelenTails
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
	p *zettel_verzeichnisse.Pool,
) (s *Store, err error) {
	s = &Store{
		lockSmith: lockSmith,
		age:       a,
		konfig:    k,
		standort:  st,
	}

	if s.hinweisen, err = hinweisen.New(st.DirZit()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if s.verzeichnisseSchwanzen, err = makeVerzeichnisseSchwanzen(k, st, s, p); err != nil {
		err = errors.Wrap(err)
		return
	}

	if s.verzeichnisseAll, err = store_verzeichnisse.MakeZettelen(k, st.DirVerzeichnisseZettelenNeue(), s, p); err != nil {
		err = errors.Wrap(err)
		return
	}

	s.indexZettelen, err = newIndexZettelen(
		st.FileVerzeichnisseZettelen(),
		s,
	)

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

	s.Transaktion.Time = ts.Now()
	s.Transaktion.Objekten = make(map[string]transaktion.Objekte)

	return
}

func (s *Store) SetZettelTransactedPrinter(ztp ZettelTransactedPrinter) {
	s.zettelTransactedPrinter = ztp
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

	f := zettel.Objekte{}

	if _, err = f.WriteTo(z, w); err != nil {
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

	// if err = s.indexZettelenTails.add(tz); err != nil {
	// 	err = errors.Wrapf(err, "failed to write zettel to index: %s", tz.Named)
	// 	return
	// }

	if err = s.verzeichnisseSchwanzen.Add(tz, tz.Named.Hinweis.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.verzeichnisseAll.Add(tz, tz.Named.Stored.Sha.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.indexZettelen.add(tz); err != nil {
		err = errors.Wrapf(err, "failed to write zettel to index: %s", tz.Named)
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

func (i *Store) ReadManySchwanzen(
	ws ...zettel_verzeichnisse.Writer,
) (err error) {
	return i.verzeichnisseSchwanzen.ReadMany(
		append(
			[]zettel_verzeichnisse.Writer{
				zettel_verzeichnisse.MakeWriterKonfig(i.konfig),
			},
			ws...,
		)...,
	)
}

func (s Store) ReadHinweisSchwanzen(
	h hinweis.Hinweis,
) (zv zettel_transacted.Zettel, err error) {
	return s.verzeichnisseSchwanzen.ReadHinweisSchwanzen(h)
}

func (s Store) ReadAllSchwanzen(ws ...zettel_transacted.Writer) (err error) {
	w := zettel_verzeichnisse.WriterZettelTransacted{
		Writer: zettel_transacted.MakeWriterMulti(
			nil,
			ws...,
		),
	}

	return s.verzeichnisseSchwanzen.Zettelen.ReadMany(w)
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

	s.zettelTransactedPrinter.ZettelTransacted(tz).Print()

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

	s.zettelTransactedPrinter.ZettelTransacted(tz).Print()

	return
}

func (s Store) Etiketten() (es []etikett.Etikett, err error) {
	return s.indexEtiketten.allEtiketten()
}

func (s *Store) Update(
	h hinweis.Hinweis,
	z zettel.Zettel,
) (tz zettel_transacted.Zettel, err error) {
	if !s.lockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "update",
		}

		return
	}

	if err = z.ApplyKonfig(s.konfig); err != nil {
		err = errors.Wrap(err)
		return
	}

	var mutter zettel_transacted.Zettel

	if mutter, err = s.verzeichnisseSchwanzen.ReadHinweisSchwanzen(h); err != nil {
		err = errors.Wrap(err)
		return
	}

	tz.Named.Hinweis = h
	tz.Named.Stored.Zettel = z

	if tz.Named.Stored.Sha, err = s.WriteZettelObjekte(z); err != nil {
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

	//TODO fix etiketten deltas
	d := etikett.MakeSetDelta(
		mutter.Named.Stored.Zettel.Etiketten,
		tz.Named.Stored.Zettel.Etiketten,
	)

	if err = s.indexEtiketten.add(d.Added); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.indexEtiketten.del(d.Removed); err != nil {
		err = errors.Wrap(err)
		return
	}

	s.zettelTransactedPrinter.ZettelTransacted(tz).Print()

	return
}

func (s Store) RevertTransaktion(t transaktion.Transaktion) (tzs zettel_transacted.Set, err error) {
	if !s.lockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "revert",
		}

		return
	}

	tzs = zettel_transacted.MakeSetUnique(len(t.Objekten))

	for _, o := range t.Objekten {
		var h *hinweis.Hinweis
		ok := false

		if h, ok = o.Id.(*hinweis.Hinweis); !ok {
			//TODO
			continue
		}

		if !o.Mutter[1].IsZero() {
			err = errors.Errorf("merges reverts are not yet supported: %s", o)
			return
		}

		errors.Print(o)

		var chain zettel_transacted.Slice

		if chain, err = s.AllInChain(*h); err != nil {
			err = errors.Wrap(err)
			return
		}

		var tz zettel_transacted.Zettel

		for _, someTz := range chain {
			errors.Print(someTz)
			if someTz.Schwanz == o.Mutter[0] {
				tz = someTz
				break
			}
		}

		if tz.Named.Stored.Sha.IsNull() {
			err = errors.Errorf("zettel not found in index!: %#v", o)
			return
		}

		if tz, err = s.Update(*h, tz.Named.Stored.Zettel); err != nil {
			err = errors.Wrap(err)
			return
		}

		tzs.Add(tz)
	}

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

	if err = s.indexZettelen.Flush(); err != nil {
		err = errors.Wrapf(err, "failed to flush new zettel index")
		return
	}

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

func (s Store) AllInChain(h hinweis.Hinweis) (c zettel_transacted.Slice, err error) {
	var mst zettel_transacted.Set

	if mst, err = s.indexZettelen.ReadHinweis(h); err != nil {
		err = errors.Wrap(err)
		return
	}

	c = mst.ToSlice()

	c.Sort(
		func(i, j int) bool { return c.Get(i).Schwanz.Less(c.Get(j).Schwanz) },
	)

	return
}

func (s *Store) ReadHinweisAt(
	h hinweis.HinweisWithIndex,
) (tz zettel_transacted.Zettel, err error) {
	if h.Index < 0 {
		errors.PrintDebug(h)
		return s.verzeichnisseSchwanzen.ReadHinweisSchwanzen(h.Hinweis)
	}

	var chain zettel_transacted.Slice

	if chain, err = s.AllInChain(h.Hinweis); err != nil {
		err = errors.Wrap(err)
		return
	}

	if chain.Len() == 0 {
		err = ErrNotFound{Id: h}
		return
	} else if chain.Len()-1 < h.Index {
		err = ErrChainIndexOutOfBounds{
			HinweisWithIndex: h,
			ChainLength:      chain.Len(),
		}

		return
	}

	tz = chain.Get(h.Index)

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

	var ts []transaktion.Transaktion

	if ts, err = s.ReadAllTransaktions(); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, t := range ts {
		errors.Print(t)

		for _, o := range t.Objekten {
			errors.Print(o)

			switch o.Gattung {

			case gattung.Zettel:
				var tz zettel_transacted.Zettel

				if tz, err = s.transactedZettelFromTransaktionObjekte(t, o); err != nil {
					if errors.Is(err, ErrNotFound{}) {
						errors.Print(err)
						err = nil
						continue
					} else {
						err = errors.Wrap(err)
						return
					}
				}

				if err = s.writeNamedZettelToIndex(tz); err != nil {
					err = errors.Wrap(err)
					return
				}

			default:
				continue
			}
		}
	}

	return
}
