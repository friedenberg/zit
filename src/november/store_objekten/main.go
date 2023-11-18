package store_objekten

import (
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/golf/kennung_index"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/erworben"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/juliett/objekte"
	"github.com/friedenberg/zit/src/kilo/objekte_store"
	"github.com/friedenberg/zit/src/kilo/store_verzeichnisse"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/mike/store_util"
	"github.com/friedenberg/zit/srx/bravo/expansion"
)

type Store struct {
	store_util.StoreUtil

	protoZettel            zettel.ProtoZettel
	verzeichnisseSchwanzen *verzeichnisseSchwanzen
	verzeichnisseAll       *store_verzeichnisse.Store
	konfigStore            konfigStore

	objekte_store.LogWriter

	isReindexing bool
	lock         sync.Locker
}

func Make(
	su store_util.StoreUtil,
) (s *Store, err error) {
	s = &Store{
		lock:      &sync.Mutex{},
		StoreUtil: su,
	}

	su.SetMatchableAdder(s)

	s.protoZettel = zettel.MakeProtoZettel(su.GetKonfig())

	if s.verzeichnisseSchwanzen, err = makeVerzeichnisseSchwanzen(
		s.StoreUtil,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if s.verzeichnisseAll, err = store_verzeichnisse.MakeStore(
		s.GetKonfig(),
		s.StoreUtil.GetStandort().DirVerzeichnisseZettelenNeue(),
		s.GetStandort(),
		nil,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	s.konfigStore.akteFormat = objekte_store.MakeAkteFormat[erworben.Akte, *erworben.Akte](
		objekte.MakeTextParserIgnoreTomlErrors[erworben.Akte](
			s.GetStandort(),
		),
		objekte.ParsedAkteTomlFormatter[erworben.Akte, *erworben.Akte]{},
		s.GetStandort(),
	)

	s.konfigStore.StoreUtil = s.StoreUtil

	errors.TodoP1("implement for other gattung")

	return
}

func (s *Store) SetLogWriter(
	lw objekte_store.LogWriter,
) {
	s.LogWriter = lw
	s.konfigStore.LogWriter = lw
}

func (s *Store) Konfig() *konfigStore {
	return &s.konfigStore
}

func (s Store) Flush() (err error) {
	if !s.GetStandort().GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: "flush",
		}

		return
	}

	if s.GetKonfig().DryRun {
		return
	}

	if s.GetKonfig().HasChanges() {
		s.verzeichnisseSchwanzen.SetNeedsFlush()
	}

	if err = s.verzeichnisseSchwanzen.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.verzeichnisseAll.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.StoreUtil.Flush(); err != nil {
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
	// TODO-P2 only set has changes if an etikett, typ, or kasten has changes
	s.GetKonfigPtr().SetHasChanges(true)

	if !s.GetStandort().GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: "update many metadatei",
		}

		return
	}

	if err = incoming.EachPtr(
		func(mwk *sku.Transacted) (err error) {
			if _, err = s.CreateOrUpdate(
				mwk,
				mwk.Kennung,
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
	gsWithHistory := gattungen.MakeMutableSet()
	gsWithoutHistory := gattungen.MakeMutableSet()

	if err = ms.GetGattungen().Each(
		func(g gattung.Gattung) (err error) {
			m, ok := ms.Get(g)

			if !ok {
				return
			}

			if m.GetSigil().IncludesHistory() {
				return gsWithHistory.Add(g)
			} else {
				return gsWithoutHistory.Add(g)
			}
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	wg := iter.MakeErrorWaitGroup()

	f1 := func(z *sku.Transacted) (err error) {
		g := gattung.Must(z.GetGattung())
		m, ok := ms.Get(g)

		if !ok {
			err = errors.Errorf("expected query to have gattung %q", g)
			return
		}

		if includeCwd && m.GetSigil().IncludesCwd() {
			var e *sku.ExternalMaybe

			if e, ok = s.GetCwdFiles().Get(z.Kennung); ok {
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

	t.Kennung = *k

	err = sku.CalculateAndSetSha(
		t,
		s.GetPersistentMetadateiFormat(),
		s.GetObjekteFormatOptions(),
	)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.addMatchableCommon(t); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.handleNew(t); err != nil {
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
	s.lock.Lock()
	defer s.lock.Unlock()

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
	t := m.GetTyp()

	if !t.IsEmpty() {
		var ti kennung_index.KennungIndex[kennung.Typ, *kennung.Typ]

		if ti, err = s.GetTypenIndex(); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = ti.StoreOne(m.GetTyp()); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = s.GetAbbrStore().AddMatchable(m); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) GetReindexFunc(
	ti kennung_index.KennungIndex[kennung.Typ, *kennung.Typ],
) func(*sku.Transacted) error {
	return func(sk *sku.Transacted) (err error) {
		errExists := s.StoreUtil.GetAbbrStore().Exists(&sk.Kennung)

		if err = s.NewOrUpdated(errExists)(sk); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = s.handleNewOrUpdatedCommit(sk, false); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = ti.StoreOne(sk.GetTyp()); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = s.GetAbbrStore().AddMatchable(sk); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}

func (s *Store) resetReindexCommon() (ti kennung_index.KennungIndex[kennung.Typ, *kennung.Typ], err error) {
	if err = s.StoreUtil.GetStandort().ResetVerzeichnisse(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if ti, err = s.GetTypenIndex(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = ti.Reset(); err != nil {
		err = errors.Wrapf(err, "failed to reset etiketten index")
		return
	}

	if err = s.StoreUtil.GetKennungIndex().Reset(); err != nil {
		err = errors.Wrapf(err, "failed to reset index kennung")
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

	s.isReindexing = true
	defer func() {
		s.isReindexing = false
	}()

	var ti kennung_index.KennungIndex[kennung.Typ, *kennung.Typ]

	if ti, err = s.resetReindexCommon(); err != nil {
		err = errors.Wrap(err)
		return
	}

	f1 := s.GetReindexFunc(ti)

	if err = s.GetBestandsaufnahmeStore().ReadAllSkus(f1); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) Reset() (err error) {
	if !s.GetStandort().GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: "reset",
		}

		return
	}

	s.isReindexing = true
	defer func() {
		s.isReindexing = false
	}()

	if _, err = s.resetReindexCommon(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
