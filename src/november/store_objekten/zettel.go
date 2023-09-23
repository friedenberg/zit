package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/todo"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/charlie/hinweisen"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/typ_akte"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/external"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/india/sku_fmt"
	"github.com/friedenberg/zit/src/juliett/objekte"
	"github.com/friedenberg/zit/src/kilo/checked_out"
	"github.com/friedenberg/zit/src/kilo/store_verzeichnisse"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/lima/objekte_store"
	"github.com/friedenberg/zit/src/mike/store_util"
)

type zettelStore struct {
	*store_util.CommonStore[
		zettel.Objekte,
		*zettel.Objekte,
		kennung.Hinweis,
		*kennung.Hinweis,
	]

	protoZettel zettel.ProtoZettel

	verzeichnisseSchwanzen *verzeichnisseSchwanzen
	verzeichnisseAll       *store_verzeichnisse.Zettelen
	tagp                   schnittstellen.AkteGetterPutter[*typ_akte.V0]
}

func makeZettelStore(
	sa store_util.StoreUtil,
	p schnittstellen.Pool[sku.Transacted2, *sku.Transacted2],
	tagp schnittstellen.AkteGetterPutter[*typ_akte.V0],
) (s *zettelStore, err error) {
	s = &zettelStore{
		protoZettel: zettel.MakeProtoZettel(sa.GetKonfig()),
		tagp:        tagp,
	}

	s.CommonStore, err = store_util.MakeCommonStore[
		zettel.Objekte,
		*zettel.Objekte,
		kennung.Hinweis,
		*kennung.Hinweis,
	](
		gattung.Zettel,
		s,
		sa,
		s,
		nil,
	)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if s.verzeichnisseSchwanzen, err = makeVerzeichnisseSchwanzen(
		s.StoreUtil,
		p,
		tagp,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if s.verzeichnisseAll, err = store_verzeichnisse.MakeZettelen(
		s.StoreUtil.GetKonfig(),
		s.StoreUtil.GetStandort().DirVerzeichnisseZettelenNeue(),
		s.StoreUtil,
		p,
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

func (s *zettelStore) AddOne(t *sku.Transacted2) (err error) {
	return s.writeNamedZettelToIndex(t)
}

func (s *zettelStore) UpdateOne(t *sku.Transacted2) (err error) {
	return s.writeNamedZettelToIndex(t)
}

func (s *zettelStore) writeNamedZettelToIndex(
	tz sku.SkuLikePtr,
) (err error) {
	errors.Log().Print("writing to index")

	if !s.StoreUtil.GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: "write named zettel to index",
		}

		return
	}

	errors.Log().Printf("writing zettel to index: %s", tz)

	s.GetKonfig().ApplyToMetadatei(tz, s.tagp)

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
) (tz *sku.Transacted2, err error) {
	var h kennung.Hinweis

	if err = h.Set(i.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	var tz1 *sku.Transacted2

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

func (i *zettelStore) ReadAllSchwanzen(
	w schnittstellen.FuncIter[sku.SkuLikePtr],
) (err error) {
	return i.verzeichnisseSchwanzen.ReadMany(
		func(z *sku.Transacted2) (err error) {
			return w(z)
		},
	)
}

func (i *zettelStore) ReadAll(
	w schnittstellen.FuncIter[sku.SkuLikePtr],
) (err error) {
	return i.verzeichnisseAll.ReadMany(
		func(z *sku.Transacted2) (err error) {
			return w(z)
		},
	)
}

func (s *zettelStore) Create(
	mg metadatei.Getter,
) (tz *sku.Transacted2, err error) {
	if !s.StoreUtil.GetLockSmith().IsAcquired() {
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

	if err = s.StoreUtil.GetKonfig().ApplyToNewMetadatei(&m, s.tagp); err != nil {
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
	incoming schnittstellen.SetLike[sku.SkuLike],
) (err error) {
	if !s.StoreUtil.GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: "update many metadatei",
		}

		return
	}

	if err = s.ReadAllSchwanzen(
		func(zt sku.SkuLikePtr) (err error) {
			ke := zt.GetKennungLike()

			if !gattung.Must(ke.GetGattung()).Equals(gattung.Zettel) {
				return
			}

			k := kennung.FormattedString(ke)

			var mwk sku.SkuLike
			ok := false

			if mwk, ok = incoming.Get(k); !ok {
				return
			}

			mwkClone := mwk.MutableClone()
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
	co objekte.ExternalLike,
) (tl sku.SkuLike, err error) {
	ze := co.(*external.Zettel)
	return s.Update(ze.GetMetadatei(), &ze.Kennung)
}

func (s *zettelStore) UpdateCheckedOut(
	co *checked_out.Zettel,
) (t *sku.Transacted2, err error) {
	errors.TodoP2("support dry run")

	if !s.StoreUtil.GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: "update",
		}

		return
	}

	m := co.External.GetMetadatei()
	m.ResetWith(m)

	if err = s.StoreUtil.GetKonfig().ApplyToNewMetadatei(&m, s.tagp); err != nil {
		err = errors.Wrap(err)
		return
	}

	co.External.SetMetadatei(m)

	if co.External.Metadatei.EqualsSansTai(co.Internal.Metadatei) {
		t = &co.Internal

		if err = s.LogWriter.Unchanged(t); err != nil {
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
) (tz *sku.Transacted2, err error) {
	errors.TodoP2("support dry run")
	var h kennung.Hinweis

	if err = h.Set(k.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !s.StoreUtil.GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: "update",
		}

		return
	}

	var mutter *sku.Transacted2

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
	mutter sku.SkuLikePtr,
) (tz *sku.Transacted2, err error) {
	if mutter == nil {
		panic("mutter was nil")
	}

	m := mg.GetMetadatei()

	if err = s.StoreUtil.GetKonfig().ApplyToNewMetadatei(&m, s.tagp); err != nil {
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

		if err = s.LogWriter.Unchanged(tz); err != nil {
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
	tz *sku.Transacted2,
	addEtikettenToIndex bool,
) (err error) {
	s.StoreUtil.CommitUpdatedTransacted(tz)

	if err = s.writeNamedZettelToIndex(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.StoreUtil.AddMatchable(tz); err != nil {
		err = errors.Wrapf(err, "failed to write zettel to index: %s", tz)
		return
	}

	if err = s.LogWriter.Updated(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *zettelStore) writeObjekte(
	mg metadatei.Getter,
	k kennung.Kennung,
) (tz *sku.Transacted2, err error) {
	if mg == nil {
		panic("metadatei.Getter was nil")
	}

	m := mg.GetMetadatei()
	m.Tai = s.StoreUtil.GetTai()

	var h kennung.Hinweis

	if err = h.Set(k.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	tz = &sku.Transacted2{
		Kennung:   kennung.Kennung2{KennungPtr: &h},
		Metadatei: m,
	}

	return
}

func (s *zettelStore) Inherit(tz *sku.Transacted2) (err error) {
	errors.Log().Printf("inheriting %s", tz)

	s.StoreUtil.CommitTransacted(tz)

	var h kennung.Hinweis

	if err = h.Set(tz.GetKennung().String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	errExists := s.StoreUtil.GetAbbrStore().Hinweis().Exists(h)

	if err = s.writeNamedZettelToIndex(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.LogWriter.NewOrUpdated(errExists)(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *zettelStore) ReindexOne(
	sk sku.SkuLike,
) (o matcher.Matchable, err error) {
	var tz *sku.Transacted2
	defer s.Pool.Put(tz)

	errors.Log().Printf("reindexing: %s", sku_fmt.String(sk))

	if tz, err = s.InflateFromSku(sk); err != nil {
		// TODO-P2 decide on how to handle format errors
		errors.Err().Print(err)
		err = nil
		// err = errors.Wrap(err)
		return
	}

	o = tz

	var h kennung.Hinweis

	if err = h.Set(tz.GetKennung().String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	errExists := s.StoreUtil.GetAbbrStore().Hinweis().Exists(h)

	if err = s.writeNamedZettelToIndex(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.StoreUtil.GetAbbrStore().AddMatchable(tz); err != nil {
		err = errors.Wrapf(err, "failed to write zettel to index: %s", tz)
		return
	}

	if err = s.LogWriter.NewOrUpdated(errExists)(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
