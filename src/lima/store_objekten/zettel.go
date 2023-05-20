package store_objekten

import (
	"io"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/checkout_mode"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/bravo/todo"
	"github.com/friedenberg/zit/src/charlie/hinweisen"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
	"github.com/friedenberg/zit/src/hotel/typ"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/kilo/cwd"
	"github.com/friedenberg/zit/src/kilo/store_util"
	"github.com/friedenberg/zit/src/kilo/store_verzeichnisse"
)

type ZettelStore interface {
	CommonStoreBase[
		zettel.Objekte,
		*zettel.Objekte,
		kennung.Hinweis,
		*kennung.Hinweis,
		zettel.Verzeichnisse,
		*zettel.Verzeichnisse,
	]

	objekte_store.Creator[*zettel.Transacted]

	objekte_store.CheckedOutUpdater[
		zettel.CheckedOut,
		*zettel.Transacted,
	]

	objekte_store.Updater[
		*kennung.Hinweis,
		*zettel.Transacted,
	]

	objekte_store.UpdaterManyMetadatei
}

type zettelStore struct {
	*commonStore[
		zettel.Objekte,
		*zettel.Objekte,
		kennung.Hinweis,
		*kennung.Hinweis,
		zettel.Verzeichnisse,
		*zettel.Verzeichnisse,
	]

	textParser  metadatei.TextParser
	protoZettel zettel.ProtoZettel

	verzeichnisseSchwanzen *verzeichnisseSchwanzen
	verzeichnisseAll       *store_verzeichnisse.Zettelen
}

func makeZettelStore(
	sa store_util.StoreUtil,
	p schnittstellen.Pool[zettel.Transacted, *zettel.Transacted],
) (s *zettelStore, err error) {
	s = &zettelStore{
		protoZettel: zettel.MakeProtoZettel(sa.GetKonfig()),
		textParser: metadatei.MakeTextParser(
			sa,
			nil, // TODO-P1 make akteFormatter
		),
	}

	s.commonStore, err = makeCommonStore[
		zettel.Objekte,
		*zettel.Objekte,
		kennung.Hinweis,
		*kennung.Hinweis,
		zettel.Verzeichnisse,
		*zettel.Verzeichnisse,
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

func (s *zettelStore) addOne(t *zettel.Transacted) (err error) {
	return s.writeNamedZettelToIndex(t)
}

func (s *zettelStore) updateOne(t *zettel.Transacted) (err error) {
	return s.writeNamedZettelToIndex(t)
}

func (s *zettelStore) writeNamedZettelToIndex(
	tz *zettel.Transacted,
) (err error) {
	errors.Log().Print("writing to index")

	if !s.StoreUtil.GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: "write named zettel to index",
		}

		return
	}

	errors.Log().Printf("writing zettel to index: %s", tz.Sku)

	if err = s.verzeichnisseSchwanzen.Add(tz, tz.Sku.Kennung.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.verzeichnisseAll.Add(tz, tz.Sku.ObjekteSha.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.StoreUtil.GetKennungIndex().AddHinweis(tz.Sku.Kennung); err != nil {
		if errors.Is(err, hinweisen.ErrDoesNotExist{}) {
			errors.Log().Printf("kennung does not contain value: %s", err)
			err = nil
		} else {
			err = errors.Wrapf(err, "failed to write zettel to index: %s", tz.Sku)
			return
		}
	}

	return
}

func (s *zettelStore) ReadOneExternal(
	e cwd.Zettel,
	t *zettel.Transacted,
) (ez zettel.External, err error) {
	var m checkout_mode.Mode

	if m, err = e.GetCheckoutMode(); err != nil {
		err = errors.Wrap(err)
		return
	}

	ez.Sku.ResetWithExternalMaybe(e)

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

	objekte.CorrectAkteShaWith(&ez, ez)

	return
}

func (s *zettelStore) readOneExternalAkte(
	ez *zettel.External,
	t *zettel.Transacted,
) (err error) {
	ez.Akte = t.Akte
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

	sh := sha.Make(aw.Sha())
	ez.SetAkteSha(sh)

	if err = s.SaveObjekte(ez); err != nil {
		err = errors.Wrapf(err, "%s", f.Name())
		return
	}

	typKonfig := s.GetKonfig().GetApproximatedTyp(
		t.GetTyp(),
	).ApproximatedOrActual()

	if typKonfig == nil {
		err = errors.Errorf("typKonfig for zettel is nil: %s", t.GetKennung())
		return
	}

	fe := typ.GetFileExtension(typKonfig)

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
	ez *zettel.External,
	t *zettel.Transacted,
) (err error) {
	var f *os.File

	if f, err = files.Open(ez.GetObjekteFD().Path); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	ez.Akte.ResetWith(t.Akte)
	ez.Metadatei.ResetWith(t.Metadatei)

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
) (tz *zettel.Transacted, err error) {
	if tz, err = s.verzeichnisseSchwanzen.ReadHinweisSchwanzen(*i); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *zettelStore) ReadAllSchwanzen(
	w schnittstellen.FuncIter[*zettel.Transacted],
) (err error) {
	return i.verzeichnisseSchwanzen.ReadMany(w)
}

func (i *zettelStore) ReadAll(
	w schnittstellen.FuncIter[*zettel.Transacted],
) (err error) {
	return i.verzeichnisseAll.ReadMany(w)
}

func (s *zettelStore) Create(
	mg metadatei.Getter,
) (tz *zettel.Transacted, err error) {
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

	if err = s.StoreUtil.GetKonfig().ApplyToMetadatei(&m); err != nil {
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

	objekte.CorrectAkteShaWith(tz, tz)

	if err = s.commitIndexMatchUpdate(tz, true); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *zettelStore) UpdateManyMetadatei(
	incoming schnittstellen.Set[metadatei.WithKennung],
) (err error) {
	if !s.StoreUtil.GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: "update many metadatei",
		}

		return
	}

	if err = s.ReadAllSchwanzen(
		func(zt *zettel.Transacted) (err error) {
			ke := zt.GetKennung()

			if !gattung.Must(ke.GetGattung()).Equals(gattung.Zettel) {
				return
			}

			k := ke.String()

			var mwk metadatei.WithKennung
			ok := false

			if mwk, ok = incoming.Get(k); !ok {
				return
			}

			mwk.Metadatei.AkteSha = sha.Make(zt.GetAkteSha())

			if _, err = s.updateLockedWithMutter(
				mwk,
				mwk.Kennung.(*kennung.Hinweis),
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
) (tl objekte.TransactedLike, err error) {
	ze := co.(*zettel.External)
	return s.Update(ze.GetMetadatei(), &ze.Sku.Kennung)
}

func (s *zettelStore) UpdateCheckedOut(
	co zettel.CheckedOut,
) (t *zettel.Transacted, err error) {
	errors.TodoP2("support dry run")

	if !s.StoreUtil.GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: "update",
		}

		return
	}

	m := co.External.GetMetadatei()
	m.ResetWith(m)

	if err = s.StoreUtil.GetKonfig().ApplyToMetadatei(&m); err != nil {
		err = errors.Wrap(err)
		return
	}

	co.External.SetMetadatei(m)

	if err = s.SaveObjekte(&co.External); err != nil {
		err = errors.Wrap(err)
		return
	}

	shaObj := sha.Make(co.External.GetObjekteSha())

	if shaObj.Equals(co.Internal.Sku.ObjekteSha) {
		t = &co.Internal

		if err = s.LogWriter.Unchanged(t); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	if t, err = s.writeObjekte(
		co.External.GetMetadatei(),
		co.External.Sku.Kennung,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	objekte.CorrectAkteShaWith(t, co.External)

	if err = s.commitIndexMatchUpdate(t, false); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *zettelStore) Update(
	mg metadatei.Getter,
	h *kennung.Hinweis,
) (tz *zettel.Transacted, err error) {
	errors.TodoP2("support dry run")

	if !s.StoreUtil.GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: "update",
		}

		return
	}

	var mutter *zettel.Transacted

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
	mutter *zettel.Transacted,
) (tz *zettel.Transacted, err error) {
	if mutter == nil {
		panic("mutter was nil")
	}

	m := mg.GetMetadatei()

	if err = s.StoreUtil.GetKonfig().ApplyToMetadatei(&m); err != nil {
		err = errors.Wrap(err)
		return
	}

	if tz, err = s.writeObjekte(m, *h); err != nil {
		err = errors.Wrap(err)
		return
	}

	objekte.CorrectAkteShaWith(tz, tz)

	shaObj := sha.Make(tz.GetObjekteSha())

	if shaObj.Equals(mutter.Sku.ObjekteSha) {
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
	tz *zettel.Transacted,
	addEtikettenToIndex bool,
) (err error) {
	s.StoreUtil.CommitUpdatedTransacted(tz)

	if err = s.writeNamedZettelToIndex(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.StoreUtil.AddMatchable(tz); err != nil {
		err = errors.Wrapf(err, "failed to write zettel to index: %s", tz.Sku)
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
) (tz *zettel.Transacted, err error) {
	if mg == nil {
		panic("metadatei.Getter was nil")
	}

	t := s.StoreUtil.GetTai()

	tz = &zettel.Transacted{
		Sku: sku.Transacted[kennung.Hinweis, *kennung.Hinweis]{
			Kennung: h,
			Kopf:    t,
			Schwanz: t,
		},
	}

	m := mg.GetMetadatei()

	tz.SetMetadatei(m)

	if err = s.SaveObjekte(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	tz.Verzeichnisse.ResetWithObjekteMetadateiGetter(tz.Akte, m)

	return
}

func (s *zettelStore) Inherit(tz *zettel.Transacted) (err error) {
	errors.Log().Printf("inheriting %s", tz.Sku.ObjekteSha)

	if err = s.SaveObjekte(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	s.StoreUtil.CommitTransacted2(tz)

	errExists := s.StoreUtil.GetAbbrStore().Hinweis().Exists(tz.Sku.Kennung)

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
	sk sku.DataIdentity,
) (o kennung.Matchable, err error) {
	var tz *zettel.Transacted
	defer s.pool.Put(tz)

	errors.Log().Printf("reindexing: %#v", o)

	if tz, err = s.InflateFromDataIdentity(sk); err != nil {
		// TODO-P2 decide on how to handle format errors
		errors.Err().Print(err)
		err = nil
		// err = errors.Wrap(err)
		return
	}

	objekte.CorrectAkteShaWith(tz, tz)

	o = tz
	errExists := s.StoreUtil.GetAbbrStore().Hinweis().Exists(tz.Sku.Kennung)

	if err = s.writeNamedZettelToIndex(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.StoreUtil.GetAbbrStore().AddMatchable(tz); err != nil {
		err = errors.Wrapf(err, "failed to write zettel to index: %s", tz.Sku)
		return
	}

	if err = s.LogWriter.NewOrUpdated(errExists)(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
