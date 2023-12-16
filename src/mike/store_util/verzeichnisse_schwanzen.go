package store_util

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/kilo/objekte_store"
	"github.com/friedenberg/zit/src/kilo/store_verzeichnisse"
)

type VerzeichnisseSchwanzen struct {
	headers [store_verzeichnisse.PageCount]*sku.Schwanzen
	*store_verzeichnisse.Store
	su StoreUtil
}

func MakeVerzeichnisseSchwanzen(
	sa StoreUtil,
) (s *VerzeichnisseSchwanzen, err error) {
	s = &VerzeichnisseSchwanzen{
		su: sa,
	}

	for i := range s.headers {
		s.headers[i] = sku.MakeSchwanzen(sa.GetKennungIndex(), s.applyKonfig)
	}

	s.Store, err = store_verzeichnisse.MakeStore(
		s.su.GetKonfig(),
		s.su.GetStandort().DirVerzeichnisseZettelenNeueSchwanzen(),
		sa.GetStandort(),
		s,
	)

	return
}

func (s *VerzeichnisseSchwanzen) Flush() (err error) {
	if err = s.Store.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.su.GetKennungIndex().Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *VerzeichnisseSchwanzen) ReadHinweisSchwanzen(
	h kennung.Kennung,
) (found *sku.Transacted, err error) {
	var n int

	if n, err = s.PageForKennung(h); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.Log().Printf("searching page %d", n)

	w := func(zv *sku.Transacted) (err error) {
		if !kennung.Equals(zv.GetKennung(), h) {
			return
		}

		found = sku.GetTransactedPool().Get()
		sku.TransactedResetter.ResetWithPtr(found, zv)

		err = collections.MakeErrStopIteration()

		return
	}

	var p *store_verzeichnisse.Page

	if p, err = s.GetPage(n); err != nil {
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

func (s *VerzeichnisseSchwanzen) applyKonfig(
	z *sku.Transacted,
) (err error) {
	if !s.su.GetKonfig().HasChanges() {
		return
	}

	s.su.GetKonfig().ApplyToSku(z)

	return
}

func (s *VerzeichnisseSchwanzen) GetVerzeichnissePageDelegate(
	n int,
) store_verzeichnisse.PageDelegate {
	return s.headers[n]
}
