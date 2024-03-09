package store_objekten

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/expansion"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/charlie/sha"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/erworben"
	"code.linenisgreat.com/zit/src/india/matcher"
	"code.linenisgreat.com/zit/src/juliett/objekte"
	"code.linenisgreat.com/zit/src/kilo/objekte_store"
	"code.linenisgreat.com/zit/src/kilo/zettel"
	"code.linenisgreat.com/zit/src/mike/store_util"
)

type Store struct {
	store_util.StoreUtil

	protoZettel      zettel.ProtoZettel
	konfigAkteFormat objekte.AkteFormat[erworben.Akte, *erworben.Akte]

	objekte_store.LogWriter
}

func Make(
	su store_util.StoreUtil,
) (s *Store, err error) {
	s = &Store{
		StoreUtil: su,
	}

	su.SetMatchableAdder(s)

	s.protoZettel = zettel.MakeProtoZettel(su.GetKonfig())

	s.konfigAkteFormat = objekte_store.MakeAkteFormat[erworben.Akte, *erworben.Akte](
		objekte.MakeTextParserIgnoreTomlErrors[erworben.Akte](
			s.GetStandort(),
		),
		objekte.ParsedAkteTomlFormatter[erworben.Akte, *erworben.Akte]{},
		s.GetStandort(),
	)

	errors.TodoP1("implement for other gattung")

	return
}

func (s *Store) SetLogWriter(lw objekte_store.LogWriter) {
	s.LogWriter = lw
}

func (s *Store) GetKonfigAkteFormat() objekte.AkteFormat[erworben.Akte, *erworben.Akte] {
	return s.konfigAkteFormat
}

func (s *Store) UpdateKonfig(
	sh schnittstellen.ShaLike,
) (kt *sku.Transacted, err error) {
	return s.CreateOrUpdateAkteSha(
		&kennung.Konfig{},
		sh,
	)
}

func (s Store) Flush(
	printerHeader schnittstellen.FuncIter[string],
) (err error) {
	if !s.GetStandort().GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: "flush",
		}

		return
	}

	if s.GetKonfig().DryRun {
		return
	}

	if err = s.StoreUtil.Flush(printerHeader); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.GetAbbrStore().Flush(); err != nil {
		errors.Err().Print(err)
		err = errors.Wrapf(err, "failed to flush abbr index")
		return
	}

	return
}

func (s *Store) UpdateManyMetadatei(
	incoming sku.TransactedSet,
) (err error) {
	// TODO-P2 [rubidium/muk "only set has changes if an etikett, typ, or kasten
	// has changes"]
	s.GetKonfig().SetHasChanges(true)

	if !s.GetStandort().GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
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
		err = objekte_store.ErrLockRequired{
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

func (s *Store) QueryWithoutCwd(
	ms matcher.Query,
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	return s.query(ms, f, false)
}

func (s *Store) QueryWithCwd(
	ms matcher.Query,
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	return s.query(ms, f, true)
}

func (s *Store) query(
	ms matcher.Query,
	f schnittstellen.FuncIter[*sku.Transacted],
	includeCwd bool,
) (err error) {
	gsWithoutHistory, gsWithHistory := ms.SplitGattungenByHistory()

	wg := iter.MakeErrorWaitGroupParallel()

	f1 := func(z *sku.Transacted) (err error) {
		g := gattung.Must(z.GetGattung())
		m, ok := ms.Get(g)

		if !ok {
			return
		}

		if includeCwd && m.GetSigil().IncludesCwd() {
			var e *sku.ExternalMaybe

			if e, ok = s.GetCwdFiles().Get(&z.Kennung); ok {
				var e2 *sku.External

				if e2, err = s.ReadOneExternal(e, z); err != nil {
					err = errors.Wrap(err)
					return
				}

				if err = z.SetFromSkuLike(&e2.Transacted); err != nil {
					err = errors.Wrap(err)
					return
				}
			}
		}

		if !m.ContainsMatchable(z) {
			return
		}

		return f(z)
	}

	wg.Do(
		func() error {
			return s.ReadAllSchwanzen(gsWithoutHistory, f1)
		},
	)

	wg.Do(
		func() error {
			return s.ReadAll(gsWithHistory, f1)
		},
	)

	return wg.GetError()
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
	errExists := s.StoreUtil.GetAbbrStore().Exists(&sk.Kennung)

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
		err = objekte_store.ErrLockRequired{
			Operation: "reindex",
		}

		return
	}

	if err = s.ResetIndexes(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.StoreUtil.GetStandort().ResetVerzeichnisse(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.StoreUtil.GetVerzeichnisse().Initialize(
		s.GetKennungIndex(),
	); err != nil {
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
