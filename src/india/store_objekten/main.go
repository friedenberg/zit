package store_objekten

import (
	"io"
	"os"
	"path"
	"reflect"
	"sort"

	"github.com/friedenberg/zit/src/alfa/logz"
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/bravo/stdprinter"
	"github.com/friedenberg/zit/src/charlie/open_file_guard"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/charlie/zk_types"
	"github.com/friedenberg/zit/src/delta/age"
	"github.com/friedenberg/zit/src/delta/etikett"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/delta/id"
	"github.com/friedenberg/zit/src/delta/ts"
	"github.com/friedenberg/zit/src/echo/age_io"
	"github.com/friedenberg/zit/src/echo/transaktion"
	"github.com/friedenberg/zit/src/echo/umwelt"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
	"github.com/friedenberg/zit/src/golf/hinweisen"
	"github.com/friedenberg/zit/src/golf/zettel_formats"
	zettel_stored "github.com/friedenberg/zit/src/golf/zettel_stored"
	"github.com/friedenberg/zit/src/hotel/collections"
)

type Store struct {
	*umwelt.Umwelt
	age.Age
	hinweisen hinweisen.Hinweisen
	*indexZettelen
	*indexZettelenTails
	*indexEtiketten
	transaktion.Transaktion
}

func (s *Store) Initialize(u *umwelt.Umwelt) (err error) {
	s.Umwelt = u

	if s.Age, err = u.Age(); err != nil {
		err = errors.Error(err)
		return
	}

	if s.hinweisen, err = hinweisen.New(s.Age, u.DirZit()); err != nil {
		err = errors.Error(err)
		return
	}

	s.indexZettelenTails, err = newIndexZettelenTails(
		u,
		u.FileVerzeichnisseZettelenSchwanzen(),
		s,
	)

	s.indexZettelen, err = newIndexZettelen(
		u,
		u.FileVerzeichnisseZettelen(),
		s,
	)

	if err != nil {
		err = errors.Wrapped(err, "failed to init zettel index")
		return
	}

	s.indexEtiketten, err = newIndexEtiketten(
		u.FileVerzeichnisseEtiketten(),
		s,
	)

	if err != nil {
		err = errors.Wrapped(err, "failed to init zettel index")
		return
	}

	s.Transaktion.Time = ts.Now()

	return
}

func (s Store) Hinweisen() hinweisen.Hinweisen {
	return s.hinweisen
}

func (s Store) WriteZettelObjekte(z zettel.Zettel) (sh sha.Sha, err error) {
	var w *age_io.Mover

	mo := age_io.MoveOptions{
		Age:                      s.Age,
		FinalPath:                s.Umwelt.DirObjektenZettelen(),
		GenerateFinalPathFromSha: true,
	}

	if w, err = age_io.NewMover(mo); err != nil {
		err = errors.Error(err)
		return
	}

	defer w.Close()

	f := zettel_formats.Objekte{}

	if _, err = f.WriteTo(z, w); err != nil {
		err = errors.Error(err)
		return
	}

	sh = w.Sha()

	return
}

func (s Store) writeNamedZettelToIndex(tz zettel_stored.Transacted) (err error) {
	logz.Printf("writing zettel to index: %s", tz.Named)

	if err = s.indexZettelenTails.add(tz); err != nil {
		err = errors.Wrapped(err, "failed to write zettel to index: %s", tz.Named)
		return
	}

	if err = s.indexZettelen.add(tz); err != nil {
		err = errors.Wrapped(err, "failed to write zettel to index: %s", tz.Named)
		return
	}

	return
}

func (s Store) Read(i id.Id) (tz zettel_stored.Transacted, err error) {
	switch tid := i.(type) {
	case sha.Sha:
		f := zettel_formats.Objekte{}

		var r io.ReadCloser

		p := id.Path(tid, s.Umwelt.DirObjektenZettelen())

		if r, err = s.ReadCloserObjekten(p); err != nil {
			err = errors.Error(err)
			return
		}

		defer stdprinter.PanicIfError(r.Close)

		if _, err = f.ReadFrom(&tz.Named.Stored.Zettel, r); err != nil {
			err = errors.Error(err)
			return
		}

	case hinweis.Hinweis:
		if tz, err = s.indexZettelenTails.Read(tid); err != nil {
			err = errors.Error(err)
			return
		}

	case hinweis.HinweisWithIndex:
		if tz, err = s.ReadHinweisAt(tid); err != nil {
			err = errors.Error(err)
			return
		}

	default:
		err = errors.Errorf("unsupported identifier: %s, %#v", i, reflect.ValueOf(i))
	}

	return
}

func (s *Store) Create(in zettel.Zettel) (tz zettel_stored.Transacted, err error) {
	if in.IsEmpty() {
		err = errors.Normalf("zettel is empty")
		return
	}

	tz.Named.Stored.Zettel = in

	if tz.Named.Stored.Sha, err = s.WriteZettelObjekte(tz.Named.Stored.Zettel); err != nil {
		err = errors.Error(err)
		return
	}

	//If the zettel exists, short circuit and return that
	// if tz2, err2 := s.Read(tz.Named.Stored.Sha); err2 == nil {
	// 	tz = tz2
	// 	return
	// }

	if tz.Named.Hinweis, err = s.hinweisen.StoreNew(tz.Named.Stored.Sha); err != nil {
		err = errors.Error(err)
		return
	}

	if tz, err = s.addZettelToTransaktion(tz.Named); err != nil {
		err = errors.Error(err)
		return
	}

	if err = s.writeNamedZettelToIndex(tz); err != nil {
		err = errors.Error(err)
		return
	}

	logz.PrintDebug(tz)

	if err = s.indexEtiketten.add(tz.Named.Stored.Zettel.Etiketten); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (s *Store) CreateWithHinweis(
	in zettel.Zettel,
	h hinweis.Hinweis,
) (tz zettel_stored.Transacted, err error) {
	if in.IsEmpty() {
		err = errors.Normalf("zettel is empty")
		return
	}

	tz.Named.Stored.Zettel = in

	if tz.Named.Stored.Sha, err = s.WriteZettelObjekte(tz.Named.Stored.Zettel); err != nil {
		err = errors.Error(err)
		return
	}

	if tz, err = s.addZettelToTransaktion(tz.Named); err != nil {
		err = errors.Error(err)
		return
	}

	if err = s.writeNamedZettelToIndex(tz); err != nil {
		err = errors.Error(err)
		return
	}

	if err = s.indexEtiketten.add(tz.Named.Stored.Zettel.Etiketten); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (s Store) Etiketten() (es []etikett.Etikett, err error) {
	return s.indexEtiketten.allEtiketten()
}

func (s *Store) Update(
	h hinweis.Hinweis,
	z zettel.Zettel,
) (tz zettel_stored.Transacted, err error) {
	var mutter zettel_stored.Transacted

	if mutter, err = s.Read(h); err != nil {
		err = errors.Error(err)
		return
	}

	tz.Named.Hinweis = h
	tz.Named.Stored.Zettel = z

	if tz.Named.Stored.Sha, err = s.WriteZettelObjekte(z); err != nil {
		err = errors.Error(err)
		return
	}

	if tz, err = s.addZettelToTransaktion(tz.Named); err != nil {
		err = errors.Error(err)
		return
	}

	if err = s.writeNamedZettelToIndex(tz); err != nil {
		err = errors.Error(err)
		return
	}

	d := mutter.Named.Stored.Zettel.Etiketten.Delta(tz.Named.Stored.Zettel.Etiketten)
	logz.Print(mutter.Named.Stored.Zettel.Etiketten)
	logz.Print(tz.Named.Stored.Zettel.Etiketten)

	if err = s.indexEtiketten.add(d.Added); err != nil {
		err = errors.Error(err)
		return
	}

	if err = s.indexEtiketten.del(d.Removed); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (s Store) Revert(h hinweis.Hinweis) (named zettel_stored.Transacted, err error) {
	return
}

func (s Store) Flush() (err error) {
	if err = s.writeTransaktion(); err != nil {
		err = errors.Wrapped(err, "failed to write transaction")
		return
	}

	if err = s.hinweisen.Flush(); err != nil {
		stdprinter.Out(err)
		err = errors.Error(err)
		return
	}

	if err = s.indexZettelenTails.Flush(); err != nil {
		err = errors.Wrapped(err, "failed to flush new zettel index")
		return
	}

	if err = s.indexZettelen.Flush(); err != nil {
		err = errors.Wrapped(err, "failed to flush new zettel index")
		return
	}

	if err = s.indexEtiketten.Flush(); err != nil {
		err = errors.Wrapped(err, "failed to flush new zettel index")
		return
	}

	return
}

func (s Store) AllInChain(h hinweis.Hinweis) (c collections.SliceTransacted, err error) {
	var mst collections.SetTransacted

	if mst, err = s.indexZettelen.ReadHinweis(h); err != nil {
		err = errors.Error(err)
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

	d := s.Umwelt.DirObjektenTransaktion()

	if headNames, err = open_file_guard.ReadDirNames(d); err != nil {
		err = errors.Error(err)
		return
	}

	for _, hn := range headNames {
		logz.Print(hn)

		var tailNames []string

		if tailNames, err = open_file_guard.ReadDirNames(d, hn); err != nil {
			err = errors.Error(err)
			return
		}

		for _, tn := range tailNames {
			logz.Print(tn)

			p := path.Join(d, hn, tn)

			var t transaktion.Transaktion

			if t, err = s.readTransaktion(p); err != nil {
				err = errors.Error(err)
				return
			}

			out = append(out, t)
		}
	}

	logz.Print("sorting")
	sort.Slice(out, func(i, j int) bool { return out[i].Time.Less(out[j].Time) })
	logz.Print("done")

	return
}

func (s *Store) ReadHinweisAt(
	h hinweis.HinweisWithIndex,
) (tz zettel_stored.Transacted, err error) {
	if h.Index < 0 {
		logz.PrintDebug(h)
		return s.indexZettelenTails.Read(h.Hinweis)
	}

	var chain collections.SliceTransacted

	if chain, err = s.AllInChain(h.Hinweis); err != nil {
		err = errors.Error(err)
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
	if err = os.RemoveAll(s.Umwelt.DirVerzeichnisse()); err != nil {
		err = errors.Wrapped(err, "failed to remove verzeichnisse dir")
		return
	}

	if err = os.MkdirAll(s.Umwelt.DirVerzeichnisse(), os.ModeDir|0755); err != nil {
		err = errors.Wrapped(err, "failed to make verzeichnisse dir")
		return
	}

	var ts []transaktion.Transaktion

	if ts, err = s.ReadAllTransaktions(); err != nil {
		err = errors.Error(err)
		return
	}

	for _, t := range ts {
		logz.Print(t)

		for _, o := range t.Objekten {
			logz.Print(o)

			switch o.Type {

			case zk_types.TypeZettel:
				var tz zettel_stored.Transacted

				if tz, err = s.transactedZettelFromTransaktionObjekte(t, o); err != nil {
					if errors.Is(err, ErrNotFound{}) {
						logz.Print(err)
						err = nil
						continue
					} else {
						err = errors.Error(err)
						return
					}
				}

				if err = s.writeNamedZettelToIndex(tz); err != nil {
					err = errors.Error(err)
					return
				}

			default:
				continue
			}
		}
	}

	if err = s.indexZettelenTails.Flush(); err != nil {
		err = errors.Wrapped(err, "failed to flush new zettel index")
		return
	}

	var tails map[hinweis.Hinweis]zettel_stored.Transacted

	if tails, err = s.ZettelenSchwanzen(); err != nil {
		err = errors.Error(err)
		return
	}

	logz.Printf("tail count: %d", len(tails))

	for _, zn := range tails {
		s.indexEtiketten.add(zn.Named.Stored.Zettel.Etiketten)
	}

	return
}
