package store_objekten

import (
	"io"
	"os"
	"path"
	"reflect"
	"sort"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/paper"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/bravo/zk_types"
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
	*indexZettelenTails
	*indexEtiketten
	*indexKennung
	*indexAbbr
	transaktion.Transaktion
}

func Make(
	lockSmith LockSmith,
	a age.Age,
	k konfig.Konfig,
	st standort.Standort,
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

	s.indexZettelenTails, err = newIndexZettelenTails(
		k,
		st.FileVerzeichnisseZettelenSchwanzen(),
		s,
	)

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

	if err = s.indexZettelenTails.add(tz); err != nil {
		err = errors.Wrapf(err, "failed to write zettel to index: %s", tz.Named)
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

func (s Store) Read(i id.Id) (tz zettel_transacted.Zettel, err error) {
	switch tid := i.(type) {
	case sha.Sha:
		f := zettel.Objekte{}

		var r io.ReadCloser

		p := id.Path(tid, s.standort.DirObjektenZettelen())

		if r, err = s.ReadCloserObjekten(p); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.PanicIfError(r.Close)

		if _, err = f.ReadFrom(&tz.Named.Stored.Zettel, r); err != nil {
			err = errors.Wrap(err)
			return
		}

	case hinweis.Hinweis:
		if tz, err = s.indexZettelenTails.Read(tid); err != nil {
			err = errors.Wrap(err)
			return
		}

	case hinweis.HinweisWithIndex:
		if tz, err = s.ReadHinweisAt(tid); err != nil {
			err = errors.Wrap(err)
			return
		}

	default:
		err = errors.Errorf("unsupported identifier: %s, %#v", i, reflect.ValueOf(i))
	}

	sh := tz.Named.Stored.Sha
	var ss string

	if ss, err = s.AbbreviateSha(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	tz.Named.Stored.Sha.Short = ss

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

	var mutter zettel_transacted.Zettel

	if mutter, err = s.Read(h); err != nil {
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

	sh := tz.Named.Stored.Sha
	var ss string

	if ss, err = s.AbbreviateSha(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	tz.Named.Stored.Sha.Short = ss

	s.zettelTransactedPrinter.ZettelTransacted(tz).Print()

	return
}

func (s Store) Revert(h hinweis.Hinweis) (named zettel_transacted.Zettel, err error) {
	if !s.lockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "revert",
		}

		return
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

	if err = s.writeTransaktion(); err != nil {
		err = errors.Wrapf(err, "failed to write transaction")
		return
	}

	if err = s.indexZettelenTails.Flush(); err != nil {
		err = errors.Wrapf(err, "failed to flush new zettel index")
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

func (s Store) ReadAllTransaktions() (out []transaktion.Transaktion, err error) {
	var headNames []string

	d := s.standort.DirObjektenTransaktion()

	if headNames, err = files.ReadDirNames(d); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, hn := range headNames {
		errors.Print(hn)

		var tailNames []string

		if tailNames, err = files.ReadDirNames(d, hn); err != nil {
			err = errors.Wrap(err)
			return
		}

		for _, tn := range tailNames {
			errors.Print(tn)

			p := path.Join(d, hn, tn)

			var t transaktion.Transaktion

			if t, err = s.readTransaktion(p); err != nil {
				err = errors.Wrap(err)
				return
			}

			out = append(out, t)
		}
	}

	errors.Print("sorting")
	sort.Slice(out, func(i, j int) bool { return out[i].Time.Less(out[j].Time) })
	errors.Print("done")

	return
}

func (s *Store) ReadHinweisAt(
	h hinweis.HinweisWithIndex,
) (tz zettel_transacted.Zettel, err error) {
	if h.Index < 0 {
		errors.PrintDebug(h)
		return s.indexZettelenTails.Read(h.Hinweis)
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

			switch o.Type {

			case zk_types.TypeZettel:
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

	if err = s.indexZettelenTails.Flush(); err != nil {
		err = errors.Wrapf(err, "failed to flush new zettel index")
		return
	}

	var tails map[hinweis.Hinweis]zettel_transacted.Zettel

	if tails, err = s.ZettelenSchwanzen(); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.Printf("tail count: %d", len(tails))

	for _, zn := range tails {
		s.indexEtiketten.add(zn.Named.Stored.Zettel.Etiketten)
	}

	return
}
