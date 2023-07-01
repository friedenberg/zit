package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/kilo/store_util"
	"github.com/friedenberg/zit/src/kilo/store_verzeichnisse"
)

type verzeichnisseSchwanzen struct {
	headers [store_verzeichnisse.PageCount]*zettel.Schwanzen
	*store_verzeichnisse.Zettelen
	su store_util.StoreUtil
}

func makeVerzeichnisseSchwanzen(
	sa store_util.StoreUtil,
	p schnittstellen.Pool[zettel.Transacted, *zettel.Transacted],
) (s *verzeichnisseSchwanzen, err error) {
	s = &verzeichnisseSchwanzen{
		su: sa,
	}

	for i := range s.headers {
		s.headers[i] = zettel.MakeSchwanzen(sa.GetKennungIndex())
	}

	s.Zettelen, err = store_verzeichnisse.MakeZettelen(
		s.su.GetKonfig(),
		s.su.GetStandort().DirVerzeichnisseZettelenNeueSchwanzen(),
		s.su,
		p,
		s,
	)

	return
}

func (s *verzeichnisseSchwanzen) Flush() (err error) {
	if err = s.Zettelen.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.su.GetKennungIndex().Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *verzeichnisseSchwanzen) ReadHinweisSchwanzen(
	h kennung.Hinweis,
) (found *zettel.Transacted, err error) {
	var n int

	if n, err = s.Zettelen.PageForHinweis(h); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.Log().Printf("searching page %d", n)

	w := func(zv *zettel.Transacted) (err error) {
		if !zv.Sku.GetKennung().Equals(h) {
			return
		}

		found = &zettel.Transacted{}
		found.ResetWith(*zv)

		err = collections.MakeErrStopIteration()

		return
	}

	var p *store_verzeichnisse.Page

	if p, err = s.Zettelen.GetPage(n); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, err = p.Copy(w); err != nil {
		err = errors.Wrap(err)
		return
	}

	if found == nil {
		err = objekte_store.ErrNotFound{Id: h}
		return
	}

	return
}

func (s *verzeichnisseSchwanzen) GetVerzeichnissePageDelegate(
	n int,
) store_verzeichnisse.PageDelegate {
	return s.headers[n]
}
