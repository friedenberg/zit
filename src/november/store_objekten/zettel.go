package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/todo"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/charlie/hinweisen"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/kilo/objekte_store"
)

func (s *Store) writeNamedZettelToIndex(
	tz *sku.Transacted,
) (err error) {
	errors.Log().Print("writing to index")

	if !s.GetStandort().GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: "write named zettel to index",
		}

		return
	}

	errors.Log().Printf("writing zettel to index: %s", tz)

	s.GetKonfigPtr().ApplyToSku(tz, s.GetAkten().GetTypV0())

	if err = s.verzeichnisseSchwanzen.AddVerzeichnisse(tz, tz.GetKennungLike().String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.verzeichnisseAll.AddVerzeichnisse(tz, tz.GetKennungLike().String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.StoreUtil.GetKennungIndex().AddHinweis(tz.GetKennungLike()); err != nil {
		if errors.Is(err, hinweisen.ErrDoesNotExist{}) {
			errors.Log().Printf("kennung does not contain value: %s", err)
			err = nil
		} else {
			err = errors.Wrapf(err, "failed to write zettel to index: %s", tz)
			return
		}
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

	var ken kennung.Hinweis

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

	if err = s.commitIndexMatchUpdate(tz); err != nil {
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

	if mutter, err = s.verzeichnisseSchwanzen.ReadHinweisSchwanzen(
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

	if err = s.commitIndexMatchUpdate(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) commitIndexMatchUpdate(
	tz *sku.Transacted,
) (err error) {
	s.CommitUpdatedTransacted(tz)

	if err = s.writeNamedZettelToIndex(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.AddMatchable(tz); err != nil {
		err = errors.Wrapf(err, "failed to write zettel to index: %s", tz)
		return
	}

	if err = s.Updated(tz); err != nil {
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
  metadatei.Resetter.ResetWithPtr(&tz.Metadatei, m)

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
