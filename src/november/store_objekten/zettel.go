package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/todo"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/charlie/hinweisen"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/kilo/objekte_store"
	"github.com/friedenberg/zit/src/kilo/store_verzeichnisse"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/mike/store_util"
)

type zettelStore struct {
	store_util.StoreUtil
	objekte_store.LogWriter

	protoZettel zettel.ProtoZettel

	verzeichnisseSchwanzen *verzeichnisseSchwanzen
	verzeichnisseAll       *store_verzeichnisse.Zettelen
}

func makeZettelStore(
	sa store_util.StoreUtil,
) (s *zettelStore, err error) {
	s = &zettelStore{
		protoZettel: zettel.MakeProtoZettel(sa.GetKonfig()),
		StoreUtil:   sa,
	}

	if s.verzeichnisseSchwanzen, err = makeVerzeichnisseSchwanzen(
		s.StoreUtil,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if s.verzeichnisseAll, err = store_verzeichnisse.MakeZettelen(
		s.GetKonfig(),
		s.StoreUtil.GetStandort().DirVerzeichnisseZettelenNeue(),
		s.GetStandort(),
		nil,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *zettelStore) Flush() (err error) {
	if err = s.verzeichnisseSchwanzen.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.verzeichnisseAll.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *zettelStore) writeNamedZettelToIndex(
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

	s.GetKonfig().ApplyToMetadatei(tz, s.GetAkten().GetTypV0())

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

func (s zettelStore) ReadOne(
	i schnittstellen.StringerGattungGetter,
) (tz *sku.Transacted, err error) {
	var h kennung.Hinweis

	if err = h.Set(i.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	var tz1 *sku.Transacted

	if tz1, err = s.verzeichnisseSchwanzen.ReadHinweisSchwanzen(h); err != nil {
		err = errors.Wrap(err)
		return
	}

	tz = sku.GetTransactedPool().Get()

	if err = tz.SetFromSkuLike(tz1); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *zettelStore) Create(
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
	s.protoZettel.Apply(&m)

	err = s.StoreUtil.GetKonfig().ApplyToNewMetadatei(
		&m,
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

	if tz, err = s.writeObjekte(
		m,
		ken,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.commitIndexMatchUpdate(tz, true); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *zettelStore) UpdateManyMetadatei(
	incoming sku.TransactedSet,
) (err error) {
	if !s.GetStandort().GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: "update many metadatei",
		}

		return
	}

	if err = s.verzeichnisseSchwanzen.ReadMany(
		func(zt *sku.Transacted) (err error) {
			ke := zt.GetKennungLike()

			if !gattung.Must(ke.GetGattung()).Equals(gattung.Zettel) {
				return
			}

			k := kennung.FormattedString(ke)

			var mwk *sku.Transacted
			ok := false

			if mwk, ok = incoming.GetPtr(k); !ok {
				return
			}

			mwkClone := sku.GetTransactedPool().Get()

			if err = mwkClone.SetFromSkuLike(mwk); err != nil {
				err = errors.Wrap(err)
				return
			}

			m := mwkClone.GetMetadateiPtr()
			m.AkteSha = sha.Make(zt.GetAkteSha())

			mwk = mwkClone

			if _, err = s.updateLockedWithMutter(
				mwk,
				ke,
				zt,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *zettelStore) updateExternal(
	ze *sku.External,
) (tl *sku.Transacted, err error) {
	return s.Update(ze.GetMetadatei(), &ze.Kennung)
}

func (s *zettelStore) UpdateCheckedOut(
	co *sku.CheckedOut,
) (t *sku.Transacted, err error) {
	errors.TodoP2("support dry run")

	if !s.GetStandort().GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: "update",
		}

		return
	}

	m := co.External.GetMetadatei()
	m.ResetWith(m)

	err = s.StoreUtil.GetKonfig().ApplyToNewMetadatei(
		&m,
		s.GetAkten().GetTypV0(),
	)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	co.External.SetMetadatei(m)

	if co.External.Metadatei.EqualsSansTai(co.Internal.Metadatei) {
		t = &co.Internal

		if err = s.Unchanged(t); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	if t, err = s.writeObjekte(
		co.External.GetMetadatei(),
		co.External.Kennung,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.commitIndexMatchUpdate(t, false); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *zettelStore) Update(
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

	if tz, err = s.updateLockedWithMutter(
		mg,
		&h,
		mutter,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *zettelStore) updateLockedWithMutter(
	mg metadatei.Getter,
	h kennung.Kennung,
	mutter *sku.Transacted,
) (tz *sku.Transacted, err error) {
	if mutter == nil {
		panic("mutter was nil")
	}

	m := mg.GetMetadatei()

	err = s.StoreUtil.GetKonfig().ApplyToNewMetadatei(
		&m,
		s.GetAkten().GetTypV0(),
	)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if tz, err = s.writeObjekte(m, h); err != nil {
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

	if err = s.commitIndexMatchUpdate(tz, false); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *zettelStore) commitIndexMatchUpdate(
	tz *sku.Transacted,
	addEtikettenToIndex bool,
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

func (s *zettelStore) writeObjekte(
	mg metadatei.Getter,
	k kennung.Kennung,
) (tz *sku.Transacted, err error) {
	if mg == nil {
		panic("metadatei.Getter was nil")
	}

	m := mg.GetMetadatei()
	m.Tai = s.GetTai()

	var h kennung.Hinweis

	if err = h.Set(k.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	tz = &sku.Transacted{
		Metadatei: m,
	}

	if err = tz.Kennung.SetWithKennung(h); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *zettelStore) Inherit(tz *sku.Transacted) (err error) {
	errors.Log().Printf("inheriting %s", tz)

	s.CommitTransacted(tz)

	var h kennung.Hinweis

	if err = h.Set(tz.GetKennung().String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	errExists := s.StoreUtil.GetAbbrStore().Hinweis().Exists(h.Parts())

	if err = s.writeNamedZettelToIndex(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.NewOrUpdated(errExists)(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *zettelStore) ReindexOne(
	tz *sku.Transacted,
) (o *sku.Transacted, err error) {
	o = tz

	var h kennung.Hinweis

	if err = h.Set(tz.GetKennung().String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	errExists := s.StoreUtil.GetAbbrStore().Hinweis().Exists(h.Parts())

	if err = s.writeNamedZettelToIndex(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.StoreUtil.GetAbbrStore().AddMatchable(tz); err != nil {
		err = errors.Wrapf(err, "failed to write zettel to index: %s", tz)
		return
	}

	if err = s.NewOrUpdated(errExists)(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
