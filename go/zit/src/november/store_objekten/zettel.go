package store_objekten

import (
	"io"

	"code.linenisgreat.com/zit-go/src/alfa/errors"
	"code.linenisgreat.com/zit-go/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit-go/src/bravo/todo"
	"code.linenisgreat.com/zit-go/src/charlie/gattung"
	"code.linenisgreat.com/zit-go/src/charlie/sha"
	"code.linenisgreat.com/zit-go/src/echo/kennung"
	"code.linenisgreat.com/zit-go/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit-go/src/hotel/sku"
	"code.linenisgreat.com/zit-go/src/kilo/objekte_store"
)

func (s *Store) CreateWithAkteString(
	mg metadatei.Getter,
  akteString string,
) (tz *sku.Transacted, err error) {
  var aw sha.WriteCloser

  if aw, err = s.GetStandort().AkteWriter(); err != nil {
    err = errors.Wrap(err)
    return
  }

  if _, err = io.WriteString(aw, akteString); err != nil {
    err = errors.Wrap(err)
    return
  }

  m := mg.GetMetadatei()
  m.SetAkteSha(aw)

  defer errors.DeferredCloser(&err, aw)

  if tz, err = s.Create(m); err != nil {
    err = errors.Wrap(err)
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

	if mg.GetMetadatei().IsEmpty() {
		err = errors.Normalf("zettel is empty")
		return
	}

	if s.protoZettel.Equals(mg.GetMetadatei()) {
		err = errors.Normalf("zettel matches protozettel")
		return
	}

	m := mg.GetMetadatei()
	s.protoZettel.Apply(m)

	if err = s.StoreUtil.GetKonfig().ApplyToNewMetadatei(
		m,
		s.GetAkten().GetTypV0(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO-P1
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

	if err = s.handleNew(
		tz,
		objekte_mode.ModeCommit,
	); err != nil {
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
