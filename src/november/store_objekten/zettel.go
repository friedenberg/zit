package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/todo"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/kilo/objekte_store"
)

func (s *Store) ReadHinweisSchwanzen(
	h kennung.Kennung,
) (found *sku.Transacted, err error) {
	var n uint8

	if n, err = s.GetVerzeichnisse().PageForKennung(h); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.Log().Printf("searching page %d", n)

	w := func(zv *sku.Transacted) (err error) {
		if !kennung.Equals(zv.GetKennung(), h) {
			return
		}

		found = sku.GetTransactedPool().Get()
		sku.TransactedResetter.ResetWith(found, zv)

		err = collections.MakeErrStopIteration()

		return
	}

	p := s.GetVerzeichnisse().GetPagePair(n)

	if err = p.Schwanzen.Copy(w); err != nil {
		err = errors.Wrap(err)
		return
	}

	if found == nil {
		err = objekte_store.ErrNotFound{Id: h}
		return
	}

	return
}

func (s *Store) Create(
	mg metadatei.Getter,
) (tz *sku.Transacted, err error) {
	if !s.GetStandort().GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: "create",
		}

		return
	}

	// if in.IsEmpty() || s.protoZettel.Equals(in) {
	// 	err = errors.Normalf("zettel is empty")
	// 	return
	// }

	m := mg.GetMetadatei()
	s.protoZettel.Apply(m)

	err = s.StoreUtil.GetKonfig().ApplyToNewMetadatei(
		m,
		s.GetAkten().GetTypV0(),
	)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	// If the zettel exists, short circuit and return that
	todo.Implement()
	// if tz2, err2 := s.ReadOne(shaObj); err2 == nil {
	// 	tz = tz2
	// 	return
	// }

	var ken *kennung.Hinweis

	if ken, err = s.StoreUtil.GetKennungIndex().CreateHinweis(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if tz, err = s.makeSku(
		m,
		ken,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.handleNew(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) Update(
	mg metadatei.Getter,
	k schnittstellen.Stringer,
) (tz *sku.Transacted, err error) {
	errors.TodoP2("support dry run")
	var h kennung.Hinweis

	if err = h.Set(k.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !s.GetStandort().GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: "update",
		}

		return
	}

	var mutter *sku.Transacted

	if mutter, err = s.ReadHinweisSchwanzen(
		h,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if mutter == nil {
		panic("mutter was nil")
	}

	m := mg.GetMetadatei()

	err = s.StoreUtil.GetKonfig().ApplyToNewMetadatei(
		m,
		s.GetAkten().GetTypV0(),
	)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if tz, err = s.makeSku(m, h); err != nil {
		err = errors.Wrap(err)
		return
	}

	mu := &mutter.Metadatei.Sha

	if err = tz.Metadatei.Mutter.SetShaLike(mu); err != nil {
		err = errors.Wrap(err)
		return
	}

	if tz.Metadatei.EqualsSansTai(mutter.GetMetadatei()) {
		if err = tz.SetFromSkuLike(mutter); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = s.Unchanged(tz); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	if err = s.handleUpdated(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) makeSku(
	mg metadatei.Getter,
	k kennung.Kennung,
) (tz *sku.Transacted, err error) {
	if mg == nil {
		panic("metadatei.Getter was nil")
	}

	m := mg.GetMetadatei()
	tz = sku.GetTransactedPool().Get()
	metadatei.Resetter.ResetWith(&tz.Metadatei, m)

	if err = tz.Kennung.SetWithKennung(k); err != nil {
		err = errors.Wrap(err)
		return
	}

	if tz.Kennung.GetGattung() != gattung.Zettel {
		err = gattung.ErrWrongType{
			ExpectedType: gattung.Zettel,
			ActualType:   gattung.Must(tz.Kennung.GetGattung()),
		}
	}

	return
}
