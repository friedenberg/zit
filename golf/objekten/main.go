package objekten

import (
	"io"
	"os"
	"path"
	"sort"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/logz"
	"github.com/friedenberg/zit/alfa/stdprinter"
	"github.com/friedenberg/zit/bravo/id"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/bravo/sha"
	"github.com/friedenberg/zit/bravo/zk_types"
	"github.com/friedenberg/zit/charlie/age"
	"github.com/friedenberg/zit/charlie/etikett"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/charlie/ts"
	"github.com/friedenberg/zit/delta/objekte"
	"github.com/friedenberg/zit/delta/umwelt"
	"github.com/friedenberg/zit/echo/transaktion"
	"github.com/friedenberg/zit/echo/zettel"
	"github.com/friedenberg/zit/foxtrot/hinweisen"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/foxtrot/zettel_formats"
	"github.com/friedenberg/zit/hotel/zettels"
)

type Store struct {
	*umwelt.Umwelt
	age.Age
	hinweisen hinweisen.Hinweisen
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
		u.FileVerzeichnisseZettelen(),
		s,
		s,
	)

	if err != nil {
		err = errors.Wrapped(err, "failed to init zettel index")
		return
	}

	s.indexEtiketten, err = newIndexEtiketten(
		u.FileVerzeichnisseEtiketten(),
		s,
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

func (s Store) writeTransaktion() (err error) {
	if len(s.Transaktion.Objekten) == 0 {
		logz.Print("not writing Transaktion as there aren't any Objekten")
		return
	}

	logz.Printf("writing Transaktion with %d Objekten", len(s.Transaktion.Objekten))

	var p string

	if p, err = id.MakeDirIfNecessary(s.Transaktion.Time, s.Umwelt.DirObjektenTransaktion()); err != nil {
		err = errors.Error(err)
		return
	}

	var w io.WriteCloser

	if w, err = s.WriteCloser(p); err != nil {
		err = errors.Error(err)
		return
	}

	defer w.Close()

	f := transaktion.Writer{Transaktion: s.Transaktion}

	if _, err = f.WriteTo(w); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (s Store) WriteZettelObjekte(z zettel.Zettel) (sh sha.Sha, err error) {
	var w *objekte.Mover

	if w, err = objekte.NewWriterMover(s.Age, s.Umwelt.DirObjektenZettelen()); err != nil {
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

func (s *Store) addZettelToTransaktion(z stored_zettel.Named) (tz stored_zettel.Transacted, err error) {
	logz.Printf("adding zettel to transaktion: %s", z.Hinweis)

	var previous stored_zettel.Transacted
	var mutter [2]ts.Time

	previous, err = s.indexZettelenTails.Read(z.Hinweis)

	if err == nil {
		mutter[0] = previous.Tail
		tz.Head = previous.Head
	} else if errors.Is(err, ErrNotFound{}) {
		err = nil
		tz.Head = s.Transaktion.Time
	} else {
		err = errors.Error(err)
		return
	}

	tz.Tail = s.Transaktion.Time
	tz.Named = z

	s.Transaktion.Objekten = append(
		s.Transaktion.Objekten,
		transaktion.Objekte{
			Type:   zk_types.TypeZettel,
			Mutter: mutter,
			Id:     &z.Hinweis,
			Sha:    z.Sha,
		},
	)

	return
}

func (s Store) writeNamedZettelToIndex(tz stored_zettel.Transacted) (err error) {
	logz.Printf("writing zettel to index: %s", tz.Named)

	if err = s.indexZettelenTails.Add(tz); err != nil {
		err = errors.Wrapped(err, "failed to write zettel to index: %s", tz.Named)
		return
	}

	return
}

func (s Store) Read(id id.Id) (tz stored_zettel.Transacted, err error) {
	switch tid := id.(type) {
	case sha.Sha:
		//TODO read from fs

	case hinweis.Hinweis:
		if tz, err = s.indexZettelenTails.Read(tid); err != nil {
			err = errors.Error(err)
			return
		}

	default:
		err = errors.Errorf("unsupported identifier: %s", id)
	}

	return
}

func (s *Store) Create(in zettel.Zettel) (tz stored_zettel.Transacted, err error) {
	if in.IsEmpty() {
		err = errors.Normal(errors.Errorf("zettel is empty"))
		return
	}

	tz.Zettel = in

	if tz.Stored.Sha, err = s.WriteZettelObjekte(tz.Zettel); err != nil {
		err = errors.Error(err)
		return
	}

	if tz.Hinweis, err = s.hinweisen.StoreNew(tz.Stored.Sha); err != nil {
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

	if err = s.indexEtiketten.Add(tz.Zettel.Etiketten); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (s *Store) CreateWithHinweis(
	in zettel.Zettel,
	h hinweis.Hinweis,
) (tz stored_zettel.Transacted, err error) {
	if in.IsEmpty() {
		err = errors.Normal(errors.Errorf("zettel is empty"))
		return
	}

	tz.Zettel = in

	if tz.Stored.Sha, err = s.WriteZettelObjekte(tz.Zettel); err != nil {
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

	if err = s.indexEtiketten.Add(tz.Zettel.Etiketten); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (s Store) Etiketten() (es []etikett.Etikett, err error) {
	return s.indexEtiketten.allEtiketten()
}

func (s Store) ZettelTails(
	qs ...stored_zettel.NamedFilter,
) (tails map[hinweis.Hinweis]stored_zettel.Transacted, err error) {
	if tails, err = s.indexZettelenTails.allTransacted(qs...); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (s *Store) Update(
	h hinweis.Hinweis,
	z zettel.Zettel,
) (tz stored_zettel.Transacted, err error) {
	var mutter stored_zettel.Transacted

	if mutter, err = s.Read(h); err != nil {
		err = errors.Error(err)
		return
	}

	tz.Hinweis = h
	tz.Zettel = z

	if tz.Sha, err = s.WriteZettelObjekte(z); err != nil {
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

	added, removed := mutter.Zettel.Etiketten.Delta(tz.Zettel.Etiketten)
	logz.Print(mutter.Zettel.Etiketten)
	logz.Print(tz.Zettel.Etiketten)

	if err = s.indexEtiketten.Add(added); err != nil {
		err = errors.Error(err)
		return
	}

	if err = s.indexEtiketten.Del(removed); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (s Store) Revert(h hinweis.Hinweis) (named stored_zettel.Transacted, err error) {
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

	if err = s.indexEtiketten.Flush(); err != nil {
		err = errors.Wrapped(err, "failed to flush new zettel index")
		return
	}

	return
}

func (s Store) AllInChain(h hinweis.Hinweis) (c zettels.Chain, err error) {
	//TODO
	// rr := indexReaderChain{
	// 	Hinweis: h,
	// }

	// if err = s.indexZettelenTails.ReadPages(&rr, h.Sha().Head()); err != nil {
	// 	err = errors.Error(err)
	// 	return
	// }

	// c.Hinweis = rr.Hinweis

	// for _, z := range rr.zettels {
	// 	c.Zettels = append(c.Zettels, z.Stored)
	// }

	return
}

func (s Store) ReadAllTransaktions() (out []transaktion.Transaktion, err error) {
	var headNames []string

	d := s.Umwelt.DirObjektenTransaktion()

	if headNames, err = open_file_guard.ReadDirNames(d); err != nil {
		err = errors.Error(err)
		return
	}

	sort.Slice(headNames, func(i, j int) bool { return headNames[i] < headNames[j] })

	for _, hn := range headNames {
		var tailNames []string

		if tailNames, err = open_file_guard.ReadDirNames(d, hn); err != nil {
			err = errors.Error(err)
			return
		}

		sort.Slice(tailNames, func(i, j int) bool { return tailNames[i] < tailNames[j] })

		for _, tn := range tailNames {
			tr := &transaktion.Reader{}
			var or io.ReadCloser

			if or, err = s.ReadCloser(path.Join(d, hn, tn)); err != nil {
				err = errors.Error(err)
				return
			}

			defer or.Close()

			if _, err = tr.ReadFrom(or); err != nil {
				err = errors.Error(err)
				return
			}

			out = append(out, tr.Transaktion)
		}
	}

	return
}

func (s *Store) Reindex() (err error) {
	if err = os.RemoveAll(s.Umwelt.DirVerzeichnisse()); err != nil {
		err = errors.Wrapped(err, "failed to remove verzeichnisse dir")
		return
	}

	if err = os.MkdirAll(s.umwelt.DirVerzeichnisse(), os.ModeDir|0755); err != nil {
		err = errors.Wrapped(err, "failed to make verzeichnisse dir")
		return
	}

	var ts []transaktion.Transaktion

	if ts, err = s.ReadAllTransaktions(); err != nil {
		err = errors.Error(err)
		return
	}

	for _, t := range ts {
		for _, o := range t.Objekten {
			switch o.Type {

			case zk_types.TypeZettel:
				var tz stored_zettel.Transacted

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

	var tails map[hinweis.Hinweis]stored_zettel.Transacted

	if tails, err = s.ZettelTails(); err != nil {
		err = errors.Error(err)
		return
	}

	logz.Printf("tail count: %d", len(tails))

	for _, zn := range tails {
		s.indexEtiketten.Add(zn.Zettel.Etiketten)
	}

	return
}
