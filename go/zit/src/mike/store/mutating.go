package store

import (
	"fmt"
	"io"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/src/bravo/expansion"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/src/charlie/files"
	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/charlie/hinweisen"
	"code.linenisgreat.com/zit/src/delta/sha"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/juliett/objekte"
	"code.linenisgreat.com/zit/src/juliett/to_merge"
)

func (s *Store) handleNewOrUpdated(
	t *sku.Transacted,
	updateType objekte_mode.Mode,
) (err error) {
	return iter.Chain(
		t,
		s.AddMatchable,
		func(t *sku.Transacted) error {
			return s.handleNewOrUpdatedCommit(t, updateType)
		},
	)
}

func (s *Store) addMutterIfNecessary(
	sk *sku.Transacted,
	ut objekte_mode.Mode,
) (err error) {
	if !sk.Metadatei.Mutter().IsNull() ||
		!ut.Contains(objekte_mode.ModeAddToBestandsaufnahme) {
		return
	}

	mutter, err := s.GetBestandsaufnahmeStore().ReadOneKennungSha(sk.GetKennung())
	defer sha.GetPool().Put(mutter)

	if err != nil {
		if collections.IsErrNotFound(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	sk.Metadatei.Mutter().ResetWith(mutter)

	return
}

func (s *Store) handleNewOrUpdatedCommit(
	t *sku.Transacted,
	mode objekte_mode.Mode,
) (err error) {
	if !s.GetStandort().GetLockSmith().IsAcquired() {
		err = objekte.ErrLockRequired{
			Operation: "write named zettel to index",
		}

		return
	}

	if mode.Contains(objekte_mode.ModeAddToBestandsaufnahme) {
		// TAI must be set before calculating objekte sha
		if mode.Contains(objekte_mode.ModeUpdateTai) {
			t.SetTai(kennung.NowTai())
		}
	}

	if err = s.addMutterIfNecessary(t, mode); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = t.CalculateObjekteShas(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if mode.Contains(objekte_mode.ModeAddToBestandsaufnahme) {
		s.CommitTransacted(t)
	}

	if mode == objekte_mode.ModeEmpty {
		if err = s.GetBestandsaufnahmeStore().WriteOneObjekteMetadatei(t); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = s.GetVerzeichnisse().ExistsOneSha(
		t.Metadatei.Sha(),
	); err == collections.ErrExists {
		return
	}

	if err = s.addToTomlIndexIfNecessary(t, mode); err != nil {
		err = errors.Wrap(err)
		return
	}

	err = nil

	g := gattung.Must(t.Kennung.GetGattung())

	switch g {
	case gattung.Konfig:
		s.GetKonfig().SetHasChanges(true)

		if err = s.GetKonfig().SetTransacted(
			t,
			s.GetAkten().GetKonfigV0(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

	case gattung.Kasten, gattung.Typ, gattung.Etikett:
		// TODO be more conservative about when konfig changes actually occurred
		s.GetKonfig().SetHasChanges(true)

		if err = s.GetKonfig().AddTransacted(t, s.GetAkten()); err != nil {
			err = errors.Wrap(err)
			return
		}

	case gattung.Zettel:
		if err = s.GetKonfig().ApplyToSku(t); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = s.kennungIndex.AddHinweis(&t.Kennung); err != nil {
			if errors.Is(err, hinweisen.ErrDoesNotExist{}) {
				errors.Log().Printf("kennung does not contain value: %s", err)
				err = nil
			} else {
				err = errors.Wrapf(err, "failed to write zettel to index: %s", t)
				return
			}
		}

	default:
		err = gattung.MakeErrUnsupportedGattung(g)
		return
	}

	if err = s.GetVerzeichnisse().Add(
		t,
		t.GetKennung().String(),
		mode,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) handleNew(
	t *sku.Transacted,
	updateType objekte_mode.Mode,
) (err error) {
	if err = s.handleNewOrUpdated(t, updateType); err != nil {
		err = errors.WrapExcept(err, collections.ErrExists)
		return
	}

	return s.New(t)
}

func (s *Store) handleUpdated(
	t *sku.Transacted,
	updateType objekte_mode.Mode,
) (err error) {
	if err = s.handleNewOrUpdated(t, updateType); err != nil {
		err = errors.WrapExcept(err, collections.ErrExists)
		return
	}

	if err = s.Updated(t); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) handleUnchanged(
	t *sku.Transacted,
) (err error) {
	return s.Unchanged(t)
}

func (s *Store) CreateOrUpdateCheckedOut(
	co *sku.CheckedOut,
) (transactedPtr *sku.Transacted, err error) {
	kennungPtr := &co.External.Kennung

	if !s.GetStandort().GetLockSmith().IsAcquired() {
		err = objekte.ErrLockRequired{
			Operation: fmt.Sprintf("create or update %s", kennungPtr),
		}

		return
	}

	transactedPtr = sku.GetTransactedPool().Get()

	if err = transactedPtr.SetFromSkuLike(&co.External); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = transactedPtr.SetAkteSha(co.External.GetAkteSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO-P2: determine why Metadatei.Etiketten can be nil
	if transactedPtr.Metadatei.EqualsSansTai(&co.Internal.Metadatei) {
		transactedPtr = &co.Internal

		if err = s.handleUnchanged(transactedPtr); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	var mutter *sku.Transacted

	if mutter, err = s.ReadOne(kennungPtr); err != nil {
		if collections.IsErrNotFound(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	} else {
		if err = transactedPtr.Metadatei.SetMutter(mutter); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = s.handleUpdated(
		transactedPtr,
		objekte_mode.ModeCommit,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) CreateOrUpdateTransacted(
	in *sku.Transacted,
) (out *sku.Transacted, err error) {
	return s.CreateOrUpdate(in, in.GetKennung())
}

// TODO-project-2022-zit-collapse_skus transition this to accepting checked out
func (s *Store) createOrUpdate(
	mg metadatei.Getter,
	kennungPtr kennung.Kennung,
	mutter *sku.Transacted,
	updateType objekte_mode.Mode,
) (transactedPtr *sku.Transacted, err error) {
	if !s.GetStandort().GetLockSmith().IsAcquired() {
		err = objekte.ErrLockRequired{
			Operation: fmt.Sprintf(
				"create or update %s",
				kennungPtr.GetGattung(),
			),
		}

		return
	}

	var m *metadatei.Metadatei

	if mg != nil {
		m = mg.GetMetadatei()
	} else {
		m = metadatei.GetPool().Get()
		defer metadatei.GetPool().Put(m)
	}

	transactedPtr = sku.GetTransactedPool().Get()
	metadatei.Resetter.ResetWith(&transactedPtr.Metadatei, m)

	if err = transactedPtr.Kennung.SetWithKennung(kennungPtr); err != nil {
		err = errors.Wrap(err)
		return
	}

	if mutter != nil {
		transactedPtr.Kopf = mutter.GetKopf()
		mu := mutter.Metadatei.Sha()

		if err = transactedPtr.Metadatei.Mutter().SetShaLike(
			mu,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		errors.TodoP4("determine if this is necessary any more")
		// transactedPtr.Sku.Kopf = s.common.GetTransaktion().Time
	}

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if mutter != nil &&
		kennung.Equals(transactedPtr.GetKennung(), mutter.GetKennung()) &&
		transactedPtr.Metadatei.EqualsSansTai(&mutter.Metadatei) {
		if err = transactedPtr.SetFromSkuLike(mutter); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = s.handleUnchanged(transactedPtr); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	if err = s.readExternalAndMergeIfNecessary(transactedPtr, mutter); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.handleNewOrUpdateWithMutter(
		transactedPtr,
		mutter,
		updateType,
	); err != nil {
		err = errors.WrapExcept(err, collections.ErrExists)
		return
	}

	return
}

func (s *Store) CreateOrUpdate(
	mg metadatei.Getter,
	kennungPtr kennung.Kennung,
) (transactedPtr *sku.Transacted, err error) {
	if !s.GetStandort().GetLockSmith().IsAcquired() {
		err = objekte.ErrLockRequired{
			Operation: fmt.Sprintf(
				"create or update %s",
				kennungPtr.GetGattung(),
			),
		}

		return
	}

	var mutter *sku.Transacted

	if mutter, err = s.ReadOne(kennungPtr); err != nil {
		if collections.IsErrNotFound(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	return s.createOrUpdate(
		mg,
		kennungPtr,
		mutter,
		objekte_mode.ModeCommit,
	)
}

func (s *Store) readExternalAndMergeIfNecessary(
	transactedPtr, mutter *sku.Transacted,
) (err error) {
	if mutter == nil {
		return
	}

	var co *sku.CheckedOut

	if co, err = s.ReadOneExternalFS(transactedPtr); err != nil {
		err = nil
		return
	}

	defer sku.GetCheckedOutPool().Put(co)

	mutterEqualsExternal := co.InternalAndExternalEqualsSansTai()

	var mode checkout_mode.Mode

	if mode, err = co.External.GetCheckoutMode(); err != nil {
		err = errors.Wrap(err)
		return
	}

	op := checkout_options.Options{
		CheckoutMode: mode,
		Force:        true,
	}

	if mutterEqualsExternal {
		if co, err = s.CheckoutOne(op, transactedPtr); err != nil {
			err = errors.Wrap(err)
			return
		}

		sku.GetCheckedOutPool().Put(co)

		return
	}

	transactedPtrCopy := sku.GetTransactedPool().Get()
	defer sku.GetTransactedPool().Put(transactedPtrCopy)

	if err = transactedPtrCopy.SetFromSkuLike(transactedPtr); err != nil {
		err = errors.Wrap(err)
		return
	}

	tm := to_merge.Sku{
		Left:   transactedPtrCopy,
		Middle: &co.Internal,
		Right:  &co.External.Transacted,
	}

	var merged sku.ExternalFDs

	merged, err = s.merge(tm)

	switch {
	case errors.Is(err, &to_merge.ErrMergeConflict{}):
		if err = tm.WriteConflictMarker(
			s.GetStandort(),
			s.GetKonfig().GetStoreVersion(),
			s.GetObjekteFormatOptions(),
			co.External.FDs.MakeConflictMarker(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

	case err != nil:
		err = errors.Wrap(err)
		return

	default:
		src := merged.Objekte.GetPath()
		dst := co.External.FDs.Objekte.GetPath()

		if err = files.Rename(src, dst); err != nil {
			return
		}
	}

	return
}

func (s *Store) CreateOrUpdateAkteSha(
	kennungPtr kennung.Kennung,
	sh schnittstellen.ShaLike,
) (transactedPtr *sku.Transacted, err error) {
	if !s.GetStandort().GetLockSmith().IsAcquired() {
		err = objekte.ErrLockRequired{
			Operation: fmt.Sprintf(
				"create or update %s",
				kennungPtr.GetGattung(),
			),
		}

		return
	}

	var mutter *sku.Transacted

	if mutter, err = s.ReadOne(kennungPtr); err != nil {
		if collections.IsErrNotFound(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	transactedPtr = sku.GetTransactedPool().Get()

	if mutter == nil {
		if err = transactedPtr.Kennung.SetWithKennung(kennungPtr); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		sku.TransactedResetter.ResetWith(transactedPtr, mutter)
	}

	transactedPtr.SetAkteSha(sh)

	return s.createOrUpdate(
		transactedPtr,
		kennungPtr,
		mutter,
		objekte_mode.ModeCommit,
	)
}

func (s *Store) handleNewOrUpdateWithMutter(
	sk, mutter *sku.Transacted,
	updateType objekte_mode.Mode,
) (err error) {
	if err = iter.Chain(
		sk,
		func(t1 *sku.Transacted) error {
			if mutter == nil {
				return s.handleNew(t1, updateType)
			} else {
				return s.handleUpdated(t1, updateType)
			}
		},
	); err != nil {
		err = errors.WrapExcept(err, collections.ErrExists)
		return
	}

	return
}

//   _____    _   _       _
//  |__  /___| |_| |_ ___| |
//    / // _ \ __| __/ _ \ |
//   / /|  __/ |_| ||  __/ |
//  /____\___|\__|\__\___|_|
//

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
		err = objekte.ErrLockRequired{
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

	if err = s.GetKonfig().ApplyToNewMetadatei(
		m,
		s.GetAkten().GetTypV0(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var ken *kennung.Hinweis

	if ken, err = s.kennungIndex.CreateHinweis(); err != nil {
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

func (s *Store) UpdateKonfig(
	sh schnittstellen.ShaLike,
) (kt *sku.Transacted, err error) {
	return s.CreateOrUpdateAkteSha(
		&kennung.Konfig{},
		sh,
	)
}

func (s *Store) UpdateManyMetadatei(
	incoming sku.TransactedSet,
) (err error) {
	if !s.GetStandort().GetLockSmith().IsAcquired() {
		err = objekte.ErrLockRequired{
			Operation: "update many metadatei",
		}

		return
	}

	if err = incoming.Each(
		func(mwk *sku.Transacted) (err error) {
			if _, err = s.CreateOrUpdateTransacted(mwk); err != nil {
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

func (s *Store) RevertTo(
	sk *sku.Transacted,
	sh *sha.Sha,
) (err error) {
	if sh.IsNull() {
		err = errors.Errorf("cannot revert to null")
		return
	}

	if !s.GetStandort().GetLockSmith().IsAcquired() {
		err = objekte.ErrLockRequired{
			Operation: "update many metadatei",
		}

		return
	}

	var mutter *sku.Transacted

	if mutter, err = s.ReadOneEnnui(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	mutter.Metadatei.Mutter().ResetWith(sk.Metadatei.Sha())

	if err = mutter.CalculateObjekteShas(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer sku.GetTransactedPool().Put(mutter)

	if _, err = s.CreateOrUpdate(
		&mutter.Metadatei,
		sk.GetKennung(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) createEtikettOrTyp(k *kennung.Kennung2) (err error) {
	switch k.GetGattung() {
	default:
		err = gattung.MakeErrUnsupportedGattung(k.GetGattung())
		return

	case gattung.Typ, gattung.Etikett:
		break
	}

	t := sku.GetTransactedPool().Get()
	defer sku.GetTransactedPool().Put(t)

	if err = t.Kennung.SetWithKennung(k); err != nil {
		err = errors.Wrap(err)
		return
	}

	// 	err = sku.CalculateAndSetSha(
	// 		t,
	// 		s.GetPersistentMetadateiFormat(),
	// 		s.GetObjekteFormatOptions(),
	// 	)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.addMatchableCommon(t); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.handleNew(t, objekte_mode.ModeCommit); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) addTyp(
	t kennung.Typ,
) (err error) {
	if err = s.GetAbbrStore().Typen().Exists(t.Parts()); err == nil {
		return
	}

	err = nil

	var k kennung.Kennung2

	if err = k.SetWithKennung(t); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.createEtikettOrTyp(&k); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) addTypAndExpanded(
	t kennung.Typ,
) (err error) {
	typenExpanded := kennung.ExpandOneSlice(&t, expansion.ExpanderRight)

	for _, t := range typenExpanded {
		if err = s.addTyp(t); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s *Store) addEtikett(
	e1 kennung.Etikett,
) (err error) {
	if e1.IsVirtual() {
		return
	}

	if err = s.GetAbbrStore().Etiketten().Exists(e1.Parts()); err == nil {
		return
	}

	err = nil

	var k kennung.Kennung2

	if err = k.SetWithKennung(e1); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.createEtikettOrTyp(&k); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) addEtikettAndExpanded(
	e kennung.Etikett,
) (err error) {
	if e.IsVirtual() {
		return
	}

	etikettenExpanded := kennung.ExpandOneSlice(&e, expansion.ExpanderRight)

	for _, e1 := range etikettenExpanded {
		if err = s.addEtikett(e1); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s *Store) addMatchableTypAndEtikettenIfNecessary(
	m *sku.Transacted,
) (err error) {
	t := m.GetTyp()

	if err = s.addTypAndExpanded(t); err != nil {
		err = errors.Wrap(err)
		return
	}

	es := iter.SortedValues[kennung.Etikett](m.Metadatei.GetEtiketten())

	for _, e := range es {
		if err = s.addEtikettAndExpanded(e); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s *Store) AddMatchable(m *sku.Transacted) (err error) {
	// s.lock.Lock()
	// defer s.lock.Unlock()

	if err = s.addMatchableTypAndEtikettenIfNecessary(m); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.addMatchableCommon(m); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) addMatchableCommon(m *sku.Transacted) (err error) {
	if err = s.AddTypToIndex(&m.Metadatei.Typ); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.GetAbbrStore().AddMatchable(m); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReindexOne(besty, sk *sku.Transacted) (err error) {
	errExists := s.GetAbbrStore().Exists(&sk.Kennung)

	if err = s.NewOrUpdated(errExists)(sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.handleNewOrUpdatedCommit(
		sk,
		objekte_mode.ModeEmpty,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.AddTypToIndex(&sk.Metadatei.Typ); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.GetAbbrStore().AddMatchable(sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// TODO-P2 add support for quiet reindexing
func (s *Store) Reindex() (err error) {
	if !s.GetStandort().GetLockSmith().IsAcquired() {
		err = objekte.ErrLockRequired{
			Operation: "reindex",
		}

		return
	}

	if err = s.ResetIndexes(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.GetStandort().ResetVerzeichnisse(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.GetVerzeichnisse().Initialize(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.GetBestandsaufnahmeStore().ReadAllSkus(
		s.ReindexOne,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
