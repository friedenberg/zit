package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/echo/hinweis"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/lima/store_verzeichnisse"
)

type verzeichnisseSchwanzen struct {
	headers [store_verzeichnisse.PageCount]*zettel.Schwanzen
	*store_verzeichnisse.Zettelen
	common *common
}

func makeVerzeichnisseSchwanzen(
	sa *common,
	p *collections.Pool[zettel.Transacted],
) (s *verzeichnisseSchwanzen, err error) {
	s = &verzeichnisseSchwanzen{
		common: sa,
	}

	for i := range s.headers {
		s.headers[i] = zettel.MakeSchwanzen()
	}

	s.Zettelen, err = store_verzeichnisse.MakeZettelen(
		s.common.Konfig(),
		s.common.Standort.DirVerzeichnisseZettelenNeueSchwanzen(),
		s.common,
		p,
		s,
	)

	return
}

func (s *verzeichnisseSchwanzen) ReadHinweisSchwanzen(
	h hinweis.Hinweis,
) (tz *zettel.Transacted, err error) {
	var n int

	if n, err = s.Zettelen.PageForHinweis(h); err != nil {
		err = errors.Wrap(err)
		return
	}

	var found *zettel.Transacted
	pool := s.Zettelen.Pool()

	w := func(zv *zettel.Transacted) (err error) {
		if !zv.Sku.Kennung.Equals(h) {
			pool.Put(zv)
			return
		}

		found = zv

		err = collections.ErrStopIteration

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

	tz = found
	tz.Objekte.Etiketten = tz.Objekte.Etiketten.Copy()

	return
}

func (s *verzeichnisseSchwanzen) ZettelTransactedWriter(
	n int,
) collections.WriterFunc[*zettel.Transacted] {
	return s.headers[n].WriteZettelTransacted
}
