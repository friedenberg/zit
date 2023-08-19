package store_objekten

import (
	"io"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/checkout_mode"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/todo"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/charlie/hinweisen"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/typ_akte"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/external"
	"github.com/friedenberg/zit/src/india/sku_formats"
	"github.com/friedenberg/zit/src/india/transacted"
	"github.com/friedenberg/zit/src/juliett/objekte"
	"github.com/friedenberg/zit/src/kilo/store_verzeichnisse"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/lima/cwd"
	"github.com/friedenberg/zit/src/lima/objekte_store"
	"github.com/friedenberg/zit/src/mike/store_util"
)

type ZettelStore interface {
	CommonStoreBase[
		zettel.Objekte,
		*zettel.Objekte,
		kennung.Hinweis,
		*kennung.Hinweis,
	]

	objekte_store.Creator[*transacted.Zettel]

	objekte_store.CheckedOutUpdater[
		*zettel.CheckedOut,
		*transacted.Zettel,
	]

	objekte_store.Updater[
		*kennung.Hinweis,
		*transacted.Zettel,
	]

	objekte_store.UpdaterManyMetadatei
}

type zettelStore struct {
	*commonStore[
		zettel.Objekte,
		*zettel.Objekte,
		kennung.Hinweis,
		*kennung.Hinweis,
	]

	textParser  metadatei.TextParser
	protoZettel zettel.ProtoZettel

	verzeichnisseSchwanzen *verzeichnisseSchwanzen
	verzeichnisseAll       *store_verzeichnisse.Zettelen
	tagp                   schnittstellen.AkteGetterPutter[*typ_akte.Akte]
}

func makeZettelStore(
	sa store_util.StoreUtil,
	p schnittstellen.Pool[transacted.Zettel, *transacted.Zettel],
	tagp schnittstellen.AkteGetterPutter[*typ_akte.Akte],
) (s *zettelStore, err error) {
	s = &zettelStore{
		protoZettel: zettel.MakeProtoZettel(sa.GetKonfig()),
		textParser: metadatei.MakeTextParser(
			sa,
			nil, // TODO-P1 make akteFormatter
		),
		tagp: tagp,
	}

	s.commonStore, err = makeCommonStore[
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

	if s.commonStore.ObjekteSaver == nil {
		panic("ObjekteSaver is nil")
	}

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if s.verzeichnisseSchwanzen, err = makeVerzeichnisseSchwanzen(
		s.StoreUtil,
		p,
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

func (s *zettelStore) addOne(t *transacted.Zettel) (err error) {
	return s.writeNamedZettelToIndex(t)
}

func (s *zettelStore) updateOne(t *transacted.Zettel) (err error) {
	return s.writeNamedZettelToIndex(t)
}

func (s *zettelStore) writeNamedZettelToIndex(
	tz *transacted.Zettel,
) (err error) {
	errors.Log().Print("writing to index")

	if !s.StoreUtil.GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: "write named zettel to index",
		}

		return
	}

	errors.Log().Printf("writing zettel to index: %s", tz)

	if err = s.verzeichnisseSchwanzen.AddVerzeichnisse(tz, tz.GetKennung().String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.verzeichnisseAll.AddVerzeichnisse(tz, tz.GetKennung().String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.StoreUtil.GetKennungIndex().AddHinweis(tz.GetKennung()); err != nil {
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

func (s *zettelStore) ReadOneExternal(
	e *cwd.Zettel,
	t *transacted.Zettel,
) (ez external.Zettel, err error) {
	var m checkout_mode.Mode

	if m, err = e.GetFDs().GetCheckoutMode(); err != nil {
		err = errors.Wrap(err)
		return
	}

	ez.ResetWithExternalMaybe(*e)

	switch m {
	case checkout_mode.ModeAkteOnly:
		if err = s.readOneExternalAkte(&ez, t); err != nil {
			err = errors.Wrap(err)
			return
		}

	case checkout_mode.ModeObjekteOnly, checkout_mode.ModeObjekteAndAkte:
		if err = s.readOneExternalObjekte(&ez, t); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s *zettelStore) readOneExternalAkte(
	ez *external.Zettel,
	t *transacted.Zettel,
) (err error) {
	ez.SetMetadatei(t.GetMetadatei())

	var aw sha.WriteCloser

	if aw, err = s.StoreUtil.AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, aw)

	var f *os.File

	if f, err = files.OpenExclusiveReadOnly(
		ez.GetAkteFD().Path,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	if _, err = io.Copy(aw, f); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh := sha.Make(aw.GetShaLike())
	ez.SetAkteSha(sh)

	if err = s.SaveObjekte(ez); err != nil {
		err = errors.Wrapf(err, "%s", f.Name())
		return
	}

	typKonfig := s.GetKonfig().GetApproximatedTyp(
		t.GetTyp(),
	).ApproximatedOrActual()

	if typKonfig == nil {
		err = errors.Errorf(
			"typKonfig for zettel is nil: %s",
			t.GetKennungLike(),
		)
		return
	}

	fe := s.GetKonfig().TypenToExtensions[t.GetTyp()]

	if fe != ez.GetAkteFD().ExtSansDot() {
		err = errors.Wrap(ErrExternalAkteExtensionMismatch{
			Expected: fe,
			Actual:   ez.GetAkteFD(),
		})

		return
	}

	return
}

func (s *zettelStore) readOneExternalObjekte(
	ez *external.Zettel,
	t *transacted.Zettel,
) (err error) {
	var f *os.File

	if f, err = files.Open(ez.GetObjekteFD().Path); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	ez.GetMetadateiPtr().ResetWith(t.GetMetadatei())

	if _, err = s.textParser.ParseMetadatei(f, ez); err != nil {
		err = errors.Wrapf(err, "%s", f.Name())
		return
	}

	if err = s.SaveObjekte(ez); err != nil {
		err = errors.Wrapf(err, "%s", f.Name())
		return
	}

	return
}

func (s zettelStore) ReadOne(
	i *kennung.Hinweis,
) (tz *transacted.Zettel, err error) {
	if tz, err = s.verzeichnisseSchwanzen.ReadHinweisSchwanzen(*i); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *zettelStore) ReadAllSchwanzen(
	w schnittstellen.FuncIter[*transacted.Zettel],
) (err error) {
	return i.verzeichnisseSchwanzen.ReadMany(w)
}

func (i *zettelStore) ReadAll(
	w schnittstellen.FuncIter[*transacted.Zettel],
) (err error) {
	return i.verzeichnisseAll.ReadMany(w)
}

func (s *zettelStore) Create(
	mg metadatei.Getter,
) (tz *transacted.Zettel, err error) {
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

	if err = s.StoreUtil.GetKonfig().ApplyToMetadatei(&m, s.tagp); err != nil {
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
		func(zt *transacted.Zettel) (err error) {
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
				mwkClone.GetKennungLikePtr().(*kennung.Hinweis),
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
	co *zettel.CheckedOut,
) (t *transacted.Zettel, err error) {
	errors.TodoP2("support dry run")

	if !s.StoreUtil.GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: "update",
		}

		return
	}

	m := co.External.GetMetadatei()
	m.ResetWith(m)

	if err = s.StoreUtil.GetKonfig().ApplyToMetadatei(&m, s.tagp); err != nil {
		err = errors.Wrap(err)
		return
	}

	co.External.SetMetadatei(m)

	if err = s.SaveObjekte(&co.External); err != nil {
		err = errors.Wrap(err)
		return
	}

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
	h *kennung.Hinweis,
) (tz *transacted.Zettel, err error) {
	errors.TodoP2("support dry run")

	if !s.StoreUtil.GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: "update",
		}

		return
	}

	var mutter *transacted.Zettel

	if mutter, err = s.verzeichnisseSchwanzen.ReadHinweisSchwanzen(
		*h,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if tz, err = s.updateLockedWithMutter(
		mg,
		h,
		mutter,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *zettelStore) updateLockedWithMutter(
	mg metadatei.Getter,
	h *kennung.Hinweis,
	mutter *transacted.Zettel,
) (tz *transacted.Zettel, err error) {
	if mutter == nil {
		panic("mutter was nil")
	}

	m := mg.GetMetadatei()

	if err = s.StoreUtil.GetKonfig().ApplyToMetadatei(&m, s.tagp); err != nil {
		err = errors.Wrap(err)
		return
	}

	if tz, err = s.writeObjekte(m, *h); err != nil {
		err = errors.Wrap(err)
		return
	}

	if tz.Metadatei.EqualsSansTai(mutter.Metadatei) {
		tz = mutter

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
	tz *transacted.Zettel,
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
	h kennung.Hinweis,
) (tz *transacted.Zettel, err error) {
	if mg == nil {
		panic("metadatei.Getter was nil")
	}

	m := mg.GetMetadatei()
	m.Tai = s.StoreUtil.GetTai()

	tz = &transacted.Zettel{
		Kennung:   h,
		Metadatei: m,
		Kopf:      m.Tai,
	}

	if err = s.SaveObjekte(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *zettelStore) Inherit(tz *transacted.Zettel) (err error) {
	errors.Log().Printf("inheriting %s", tz)

	if err = s.SaveObjekte(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	s.StoreUtil.CommitTransacted(tz)

	errExists := s.StoreUtil.GetAbbrStore().Hinweis().Exists(
		tz.GetKennung(),
	)

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
) (o kennung.Matchable, err error) {
	var tz *transacted.Zettel
	defer s.pool.Put(tz)

	errors.Log().Printf("reindexing: %s", sku_formats.String(sk))

	if tz, err = s.InflateFromSku(sk); err != nil {
		// TODO-P2 decide on how to handle format errors
		errors.Err().Print(err)
		err = nil
		// err = errors.Wrap(err)
		return
	}

	o = tz

	errExists := s.StoreUtil.GetAbbrStore().Hinweis().Exists(
		tz.GetKennung(),
	)

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
