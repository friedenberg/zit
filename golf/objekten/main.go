package objekten

import (
	"os"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/bravo/sha"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/delta/objekte"
	"github.com/friedenberg/zit/delta/umwelt"
	"github.com/friedenberg/zit/echo/zettel"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/foxtrot/zettel_formats"
	"github.com/friedenberg/zit/hotel/zettels"
	"github.com/friedenberg/zit/verzeichnisse"
)

type Store struct {
	umwelt *umwelt.Umwelt
	zettels.Zettels
	zettelIndex *verzeichnisse.Index
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

	return
}

func (s Store) writeZettelObjekte(z zettel.Zettel) (err error) {
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

func (s Store) Create(in zettel.Zettel) (z stored_zettel.Named, err error) {
	if z, err = s.Zettels.Create(in); err != nil {
		err = errors.Error(err)
		return
	}

	if err = s.writeZettelObjekte(in); err != nil {
		err = errors.Error(err)
		return
	}

	if err = s.writeNamedZettelToIndex(z); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (s Store) CreateWithHinweis(in zettel.Zettel, h hinweis.Hinweis) (z stored_zettel.Named, err error) {
	if z, err = s.Zettels.CreateWithHinweis(in, h); err != nil {
		err = errors.Error(err)
		return
	}

	if err = s.writeZettelObjekte(in); err != nil {
		err = errors.Error(err)
		return
	}

	if err = s.writeNamedZettelToIndex(z); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (s Store) Update(z stored_zettel.Named) (stored stored_zettel.Named, err error) {
	if stored, err = s.Zettels.Update(z); err != nil {
		err = errors.Error(err)
		return
	}

	if err = s.writeZettelObjekte(stored.Zettel); err != nil {
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
	if err = s.Zettels.Flush(); err != nil {
		err = errors.Wrapped(err, "failed to flush old zettel store")
		return
	}

	if s.zettelIndex == nil {
		return
	}

	if err = s.zettelIndex.Flush(); err != nil {
		err = errors.Wrapped(err, "failed to flush new zettel index")
		return
	}

	return
}

func (s Store) Reindex() (err error) {
	if err = os.RemoveAll(s.umwelt.DirVerzeichnisseZettelen()); err != nil {
		err = errors.Wrapped(err, "failed to remove zettel index")
		return
	}

	if err = os.MkdirAll(s.umwelt.DirVerzeichnisseZettelen(), os.ModeDir|0755); err != nil {
		err = errors.Wrapped(err, "failed to make zettel index dir")
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

			if err = s.writeNamedZettelToIndex(nz); err != nil {
				err = errors.Error(err)
				return
			}
		}
	}

	return
}
