package store_objekten

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/india/zettel"
	"github.com/friedenberg/zit/src/juliett/zettel_verzeichnisse"
	"github.com/friedenberg/zit/src/kilo/store_verzeichnisse"
)

type verzeichnisseSchwanzen struct {
	headers [store_verzeichnisse.PageCount]*zettel_verzeichnisse.WriterSchwanzen
	*store_verzeichnisse.Zettelen
	common *common
}

func makeVerzeichnisseSchwanzen(
	sa *common,
	p zettel_verzeichnisse.Pool,
) (s *verzeichnisseSchwanzen, err error) {
	s = &verzeichnisseSchwanzen{
		common: sa,
	}

	for i, _ := range s.headers {
		s.headers[i] = zettel_verzeichnisse.MakeWriterSchwanzen()
	}

	s.Zettelen, err = store_verzeichnisse.MakeZettelen(
		s.common.Konfig,
		s.common.Standort.DirVerzeichnisseZettelenNeueSchwanzen(),
		s.common,
		p,
		s,
	)

	return
}

func (s *verzeichnisseSchwanzen) ReadHinweisSchwanzen(
	h hinweis.Hinweis,
) (tz zettel.Transacted, err error) {
	var n int

	if n, err = s.Zettelen.PageForHinweis(h); err != nil {
		err = errors.Wrap(err)
		return
	}

	var found *zettel_verzeichnisse.Zettel
	pool := s.Zettelen.Pool()

	w := func(zv *zettel_verzeichnisse.Zettel) (err error) {
		if !zv.Transacted.Sku.Kennung.Equals(&h) {
			pool.Put(zv)
			return
		}

		found = zv

		err = io.EOF

		return
	}

	var p *store_verzeichnisse.Page

	if p, err = s.Zettelen.GetPage(n); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = p.WriteZettelenTo(w); err != nil {
		err = errors.Wrap(err)
		return
	}

	if found == nil {
		err = ErrNotFound{Id: h}
		return
	}

	tz = found.Transacted
	tz.Objekte.Etiketten = tz.Objekte.Etiketten.Copy()

	return
}

func (s *verzeichnisseSchwanzen) ZettelVerzeichnisseWriter(
	n int,
) collections.WriterFunc[*zettel_verzeichnisse.Zettel] {
	return s.headers[n].WriteZettelVerzeichnisse
}
