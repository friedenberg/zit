package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/typ_akte"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/india/transacted"
	"github.com/friedenberg/zit/src/kilo/store_verzeichnisse"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/lima/objekte_store"
	"github.com/friedenberg/zit/src/mike/store_util"
)

type verzeichnisseSchwanzen struct {
	headers [store_verzeichnisse.PageCount]*zettel.Schwanzen
	tagp    schnittstellen.AkteGetterPutter[*typ_akte.V0]
	*store_verzeichnisse.Zettelen
	su store_util.StoreUtil
}

func makeVerzeichnisseSchwanzen(
	sa store_util.StoreUtil,
	p schnittstellen.Pool[transacted.Zettel, *transacted.Zettel],
	tagp schnittstellen.AkteGetterPutter[*typ_akte.V0],
) (s *verzeichnisseSchwanzen, err error) {
	s = &verzeichnisseSchwanzen{
		su:   sa,
		tagp: tagp,
	}

	for i := range s.headers {
		s.headers[i] = zettel.MakeSchwanzen(sa.GetKennungIndex(), s.applyKonfig)
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
) (found *transacted.Zettel, err error) {
	var n int

	if n, err = s.Zettelen.PageForHinweis(h); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.Log().Printf("searching page %d", n)

	w := func(zv *transacted.Zettel) (err error) {
		if !zv.GetKennung().Equals(h) {
			return
		}

		found = &transacted.Zettel{}
		found.ResetWith(*zv)

		err = collections.MakeErrStopIteration()

		return
	}

	var p *store_verzeichnisse.Page

	if p, err = s.Zettelen.GetPage(n); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = p.Copy(w); err != nil {
		err = errors.Wrap(err)
		return
	}

	if found == nil {
		err = objekte_store.ErrNotFound{Id: h}
		return
	}

	return
}

func (s *verzeichnisseSchwanzen) applyKonfig(
	z *transacted.Zettel,
) (err error) {
	if !s.su.GetKonfig().HasChanges() {
		return
	}

	s.su.GetKonfig().ApplyToMetadatei(z, s.tagp)

	return
}

func (s *verzeichnisseSchwanzen) GetVerzeichnissePageDelegate(
	n int,
) store_verzeichnisse.PageDelegate {
	return s.headers[n]
}
