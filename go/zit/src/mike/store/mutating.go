package store

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/expansion"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/delta/file_lock"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/delta/hinweisen"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/store_fs"
)

func (s *Store) tryRealize(
	kinder, mutter *sku.Transacted,
	o ObjekteOptions,
) (err error) {
	if mutter == nil && o.Contains(objekte_mode.ModeApplyProto) {
		s.protoZettel.Apply(kinder, kinder)
	}

	if err = s.tryPreCommitHooks(kinder, mutter, o); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO just just mutter == nil
	if mutter == nil {
		if err = s.tryNewHook(kinder, o); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = kinder.CalculateObjekteShas(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// TODO add RealizeAndOrStore result
func (s *Store) tryRealizeAndOrStore(
	kinder *sku.Transacted,
	o ObjekteOptions,
) (err error) {
	ui.Log().Printf("%s -> %s", o, kinder)

	if kinder.Kennung.IsEmpty() &&
		o.ContainsAny(
			objekte_mode.ModeAddToBestandsaufnahme,
		) {
		err = errors.Errorf("empty kennung")
		return
	}

	if !s.GetStandort().GetLockSmith().IsAcquired() &&
		o.ContainsAny(
			objekte_mode.ModeAddToBestandsaufnahme,
		) {
		err = errors.Wrap(file_lock.ErrLockRequired{
			Operation: "commit",
		})

		return
	}

	// TAI must be set before calculating objekte sha
	if o.ContainsAny(
		objekte_mode.ModeAddToBestandsaufnahme,
		objekte_mode.ModeUpdateTai,
	) {
		if o.Clock == nil {
			o.Clock = s
		}

		kinder.SetTai(o.Clock.GetTai())
	}

	var mutter *sku.Transacted

	if mutter, err = s.fetchMutterIfNecessary(kinder, o); err != nil {
		err = errors.Wrap(err)
		return
	}

	if mutter != nil {
		defer sku.GetTransactedPool().Put(mutter)
	}

	if err = s.tryRealize(kinder, mutter, o); err != nil {
		err = errors.Wrap(err)
		return
	}

	if o.Contains(objekte_mode.ModeAddToBestandsaufnahme) {
		if err = s.addMatchableTypAndEtikettenIfNecessary(kinder); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = s.addMatchableCommon(kinder); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	for _, vs := range s.virtualStores {
		if err = vs.CommitTransacted(kinder, mutter); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if o.Mode != objekte_mode.ModeReindex &&
		mutter != nil &&
		kennung.Equals(kinder.GetKennung(), mutter.GetKennung()) &&
		kinder.Metadatei.EqualsSansTai(&mutter.Metadatei) {

		if err = kinder.SetFromSkuLike(mutter); err != nil {
			err = errors.Wrap(err)
			return
		}

		if o.Mode.Contains(objekte_mode.ModeSchwanz) {
			if err = s.Unchanged(kinder); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		return
	}

	if err = s.GetKonfig().ApplySchlummerndAndRealizeEtiketten(
		kinder,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if o.Mode.Contains(objekte_mode.ModeSchwanz) {
		if err = s.GetKonfig().AddTransacted(
			kinder,
			mutter,
			s.GetAkten(),
			o.Mode,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		if kinder.GetGattung() == gattung.Zettel {
			if err = s.kennungIndex.AddHinweis(&kinder.Kennung); err != nil {
				if errors.Is(err, hinweisen.ErrDoesNotExist{}) {
					ui.Log().Printf("kennung does not contain value: %s", err)
					err = nil
				} else {
					err = errors.Wrapf(err, "failed to write zettel to index: %s", kinder)
					return
				}
			}
		}

	}

	if o.Contains(objekte_mode.ModeAddToBestandsaufnahme) {
		ui.Log().Print("adding to bestandsaufnahme", o, kinder)
		if err = s.commitTransacted(kinder, mutter); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = s.addToTomlIndexIfNecessary(kinder, o); err != nil {
			err = errors.Wrap(err)
			return
		}

		err = nil
	}

	if o.Contains(objekte_mode.ModeSchwanz) {
		if err = s.GetVerzeichnisse().Add(
			kinder,
			kinder.GetKennung().String(),
			o.Mode,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		if mutter == nil {
			if err = s.New(kinder); err != nil {
				err = errors.Wrap(err)
				return
			}
		} else {
			if err = s.Updated(kinder); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

	}

	if o.Contains(objekte_mode.ModeMergeCheckedOut) {
		if err = s.readExternalAndMergeIfNecessary(kinder, mutter); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s *Store) fetchMutterIfNecessary(
	sk *sku.Transacted,
	ut ObjekteOptions,
) (mutter *sku.Transacted, err error) {
	if !sk.Metadatei.Mutter().IsNull() {
		return
	}

	if mutter, err = s.GetVerzeichnisse().ReadOneKennung(
		sk.GetKennung(),
	); err != nil {
		if collections.IsErrNotFound(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	for _, vs := range s.virtualStores {
		if err = vs.ModifySku(mutter); err != nil {
			ui.Err().Print(err)
			err = nil
			return
		}
	}

	sk.Metadatei.Mutter().ResetWith(mutter.Metadatei.Sha())

	return
}

// TODO add results for which stores had which change types
func (s *Store) commitTransacted(
	kinder *sku.Transacted,
	mutter *sku.Transacted,
) (err error) {
	sk := sku.GetTransactedPool().Get()

	if err = sk.SetFromSkuLike(kinder); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.bestandsaufnahmeAkte.Skus.Add(sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) AddTypToIndex(t *kennung.Typ) (err error) {
	if t == nil {
		return
	}

	if t.IsEmpty() {
		return
	}

	if err = s.typenIndex.StoreOne(*t); err != nil {
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
	co *store_fs.CheckedOut,
	updateCheckout bool,
) (transactedPtr *sku.Transacted, err error) {
	kennungPtr := &co.External.Kennung

	if !s.GetStandort().GetLockSmith().IsAcquired() {
		err = file_lock.ErrLockRequired{
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

	if err = s.tryRealizeAndOrStore(
		transactedPtr,
		ObjekteOptions{Mode: objekte_mode.ModeCreate},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !updateCheckout {
		return
	}

	var mode checkout_mode.Mode

	if mode, err = co.External.GetCheckoutMode(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, err = s.CheckoutOneFS(
		checkout_options.Options{CheckoutMode: mode, Force: true},
		transactedPtr,
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

	if err = s.tryRealizeAndOrStore(
		t,
		ObjekteOptions{Mode: objekte_mode.ModeCommit},
	); err != nil {
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

func (s *Store) addEtikettAndExpanded(
	e kennung.Etikett,
) (err error) {
	if e.IsVirtual() {
		return
	}

	etikettenExpanded := kennung.ExpandOneSlice(&e, expansion.ExpanderRight)

	s.etikettenLock.Lock()
	defer s.etikettenLock.Unlock()

	for _, e1 := range etikettenExpanded {
		if e1.IsVirtual() {
			continue
		}

		if err = s.GetAbbrStore().Etiketten().Exists(e1.Parts()); err == nil {
			continue
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

	es := iter.SortedValues(m.Metadatei.GetEtiketten())

	for _, e := range es {
		if err = s.addEtikettAndExpanded(e); err != nil {
			err = errors.Wrap(err)
			return
		}
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

func (s *Store) reindexOne(besty, sk *sku.Transacted) (err error) {
	o := ObjekteOptions{
		Mode: objekte_mode.ModeReindex,
	}

	if err = s.tryRealizeAndOrStore(
		sk,
		o,
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
