package objekten

import (
	"io"
	"os"
	"path"
	"sort"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/logz"
	"github.com/friedenberg/zit/bravo/id"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/bravo/sha"
	"github.com/friedenberg/zit/bravo/zk_types"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/delta/objekte"
	"github.com/friedenberg/zit/delta/umwelt"
	"github.com/friedenberg/zit/echo/transaktion"
	"github.com/friedenberg/zit/echo/zettel"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/foxtrot/zettel_formats"
	"github.com/friedenberg/zit/hotel/zettels"
	"github.com/friedenberg/zit/ts"
	"github.com/friedenberg/zit/verzeichnisse"
)

type Store struct {
	umwelt *umwelt.Umwelt
	zettels.Zettels
	zettelIndex *verzeichnisse.Index
	transaktion.Transaktion
}

func (s *Store) Initialize(u *umwelt.Umwelt) (err error) {
	s.umwelt = u

	s.zettelIndex, err = verzeichnisse.NewIndex(
		u.DirVerzeichnisseZettelen(),
		s,
		s,
		sha.Head,
	)

	if err != nil {
		err = errors.Wrapped(err, "failed to init zettel index")
		return
	}

	s.Transaktion.Time = ts.Now()

	return
}

func (s Store) writeTransaktion() (err error) {
	if len(s.Transaktion.Objekten) == 0 {
		logz.Print("not writing Transaktion as there aren't any Objekten")
		return
	}

	logz.Printf("writing Transaktion with %d Objekten", len(s.Transaktion.Objekten))

	var p string

	if p, err = id.MakeDirIfNecessary(s.Transaktion.Time, s.Umwelt().DirObjektenTransaktion()); err != nil {
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

func (s Store) writeZettelObjekte(z zettel.Zettel) (sh sha.Sha, err error) {
	var w *objekte.Mover

	if w, err = objekte.NewWriterMover(s.Age(), s.Umwelt().DirObjektenZettelen()); err != nil {
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

func (s *Store) addZettelToTransaktion(z stored_zettel.Named) (err error) {
	logz.Printf("adding zettel to transaktion: %s", z.Hinweis)

	s.Transaktion.Objekten = append(
		s.Transaktion.Objekten,
		transaktion.Objekte{
			Type: zk_types.TypeZettel,
			Id:   &z.Hinweis,
			Sha:  z.Sha,
		},
	)

	return
}

func (s Store) writeNamedZettelToIndex(nz stored_zettel.Named) (err error) {
	rowMaker := func() ([]verzeichnisse.Row, error) {
		return s.indexRowsForZettel(nz)
	}

	if err = s.zettelIndex.WriteRows(rowMaker); err != nil {
		err = errors.Wrapped(err, "failed to write index rows for zettel %s", nz)
		return
	}

	return
}

func (s Store) Read(id id.Id) (sz stored_zettel.Named, err error) {
	rr := indexReaderOneZettel{
		Id: id,
	}

	if err = s.zettelIndex.ReadPages(&rr, id.Sha().Head()); err != nil {
		err = errors.Error(err)
		return
	}

	sz = rr.Named

	return
}

func (s *Store) Create(in zettel.Zettel) (z stored_zettel.Named, err error) {
	if z, err = s.Zettels.Create(in); err != nil {
		err = errors.Error(err)
		return
	}

	if z.Stored.Sha, err = s.writeZettelObjekte(in); err != nil {
		err = errors.Error(err)
		return
	}

	if err = s.addZettelToTransaktion(z); err != nil {
		err = errors.Error(err)
		return
	}

	if err = s.writeNamedZettelToIndex(z); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (s *Store) CreateWithHinweis(in zettel.Zettel, h hinweis.Hinweis) (z stored_zettel.Named, err error) {
	if z, err = s.Zettels.CreateWithHinweis(in, h); err != nil {
		err = errors.Error(err)
		return
	}

	if z.Stored.Sha, err = s.writeZettelObjekte(in); err != nil {
		err = errors.Error(err)
		return
	}

	if err = s.addZettelToTransaktion(z); err != nil {
		err = errors.Error(err)
		return
	}

	if err = s.writeNamedZettelToIndex(z); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (s *Store) Update(z stored_zettel.Named) (stored stored_zettel.Named, err error) {
	if stored, err = s.Zettels.Update(z); err != nil {
		err = errors.Error(err)
		return
	}

	if stored.Sha, err = s.writeZettelObjekte(stored.Zettel); err != nil {
		err = errors.Error(err)
		return
	}

	if err = s.addZettelToTransaktion(stored); err != nil {
		err = errors.Error(err)
		return
	}

	if err = s.writeNamedZettelToIndex(stored); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (s Store) Flush() (err error) {
	if err = s.writeTransaktion(); err != nil {
		err = errors.Wrapped(err, "failed to write transaction")
		return
	}

	if err = s.Zettels.Flush(); err != nil {
		err = errors.Wrapped(err, "failed to flush old zettel store")
		return
	}

	if s.zettelIndex == nil {
		err = errors.Errorf("index was not initialized")
		return
	}

	if err = s.zettelIndex.Flush(); err != nil {
		err = errors.Wrapped(err, "failed to flush new zettel index")
		return
	}

	return
}

func (s Store) AllInChain(h hinweis.Hinweis) (c zettels.Chain, err error) {
	rr := indexReaderChain{
		Hinweis: h,
	}

	if err = s.zettelIndex.ReadPages(&rr, h.Sha().Head()); err != nil {
		err = errors.Error(err)
		return
	}

	c.Hinweis = rr.Hinweis
	c.Zettels = rr.zettels

	return

}

func (s Store) ReadAllTransaktions() (out []transaktion.Transaktion, err error) {
	var headNames []string

	d := s.Umwelt().DirObjektenTransaktion()

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
	if err = os.RemoveAll(s.Umwelt().DirVerzeichnisseZettelen()); err != nil {
		err = errors.Wrapped(err, "failed to remove zettel index")
		return
	}

	if err = os.MkdirAll(s.umwelt.DirVerzeichnisseZettelen(), os.ModeDir|0755); err != nil {
		err = errors.Wrapped(err, "failed to make zettel index dir")
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

				if err = s.writeNamedZettelToIndex(tz.Named); err != nil {
					err = errors.Error(err)
					return
				}

			default:
				continue
			}
		}
	}

	return
}

func (s *Store) Rewrite() (err error) {
	if err = os.RemoveAll(s.Umwelt().DirObjektenTransaktion()); err != nil {
		err = errors.Wrapped(err, "failed to remove transaktion dir")
		return
	}

	if err = os.MkdirAll(s.umwelt.DirObjektenTransaktion(), os.ModeDir|0755); err != nil {
		err = errors.Wrapped(err, "failed to make transaktion dir")
		return
	}

	if err = os.RemoveAll(s.Umwelt().DirObjektenZettelen()); err != nil {
		err = errors.Wrapped(err, "failed to remove zettelen dir")
		return
	}

	if err = os.MkdirAll(s.umwelt.DirObjektenZettelen(), os.ModeDir|0755); err != nil {
		err = errors.Wrapped(err, "failed to make zettelen dir")
		return
	}

	var hins []hinweis.Hinweis

	if _, hins, err = s.Zettels.Hinweisen().All(); err != nil {
		err = errors.Error(err)
		return
	}

	for _, h := range hins {
		var chain zettels.Chain

		if chain, err = s.Zettels.AllInChain(h); err != nil {
			if errors.Is(err, zettels.ErrShaNotFound{}) {
				err = nil
			} else if errors.Is(err, zettels.ErrHistoryLoopDetected{}) {
				err = nil
			} else {
				err = errors.Error(err)
				return
			}
		}

		for i := len(chain.Zettels) - 1; i >= 0; i-- {
			sz := chain.Zettels[i]

			nz := stored_zettel.Named{
				Hinweis: h,
				Stored:  sz,
			}

			if nz.Stored.Sha, err = s.writeZettelObjekte(sz.Zettel); err != nil {
				err = errors.Error(err)
				return
			}

			if err = s.addZettelToTransaktion(nz); err != nil {
				err = errors.Error(err)
				return
			}

			if err = s.writeNamedZettelToIndex(nz); err != nil {
				err = errors.Error(err)
				return
			}
		}
	}

	return
}
