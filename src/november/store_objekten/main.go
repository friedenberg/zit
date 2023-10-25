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
	"github.com/friedenberg/zit/src/mike/store_util"
	"github.com/friedenberg/zit/srx/bravo/expansion"
)

type Store struct {
	store_util.StoreUtil

	zettelStore  *zettelStore
	typStore     typStore
	etikettStore etikettStore
	konfigStore  konfigStore
	kastenStore  kastenStore

	objekte_store.LogWriter

	// Gattungen
	flushers map[schnittstellen.GattungLike]errors.Flusher

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

	if s.zettelStore, err = makeZettelStore(s.StoreUtil); err != nil {
		err = errors.Wrap(err)
		return
	}

	s.konfigStore = konfigStore{
		akteFormat: objekte_store.MakeAkteFormat[erworben.Akte, *erworben.Akte](
			objekte.MakeTextParserIgnoreTomlErrors[erworben.Akte](
				s.GetStandort(),
			),
			objekte.ParsedAkteTomlFormatter[erworben.Akte, *erworben.Akte]{},
			s.GetStandort(),
		),
		StoreUtil: s.StoreUtil,
	}

	s.typStore.StoreUtil = s.StoreUtil
	s.etikettStore.StoreUtil = s.StoreUtil
	s.kastenStore.StoreUtil = s.StoreUtil

	errors.TodoP1("implement for other gattung")

	s.flushers = map[schnittstellen.GattungLike]errors.Flusher{
		gattung.Zettel:  s.zettelStore,
		gattung.Typ:     s.typStore,
		gattung.Etikett: s.etikettStore,
		gattung.Kasten:  s.kastenStore,
		gattung.Konfig:  s.konfigStore,
	}

	return
}

func (s *Store) SetLogWriter(
	lw objekte_store.LogWriter,
) {
	s.LogWriter = lw
	s.zettelStore.LogWriter = lw
	s.konfigStore.LogWriter = lw
}

func (s *Store) Zettel() *zettelStore {
	return s.zettelStore
}

func (s *Store) Typ() *typStore {
	return &s.typStore
}

func (s *Store) Etikett() *etikettStore {
	return &s.etikettStore
}

func (s *Store) Konfig() *konfigStore {
	return &s.konfigStore
}

func (s *Store) Kasten() *kastenStore {
	return &s.kastenStore
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

	for _, fl := range s.flushers {
		if err = fl.Flush(); err != nil {
			err = errors.Wrap(err)
			return
		}
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

func (s *Store) Query(
	ms matcher.Query,
	f schnittstellen.FuncIter[*sku.Transacted],
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

	if err = s.onNew(t); err != nil {
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

	es := iter.SortedValues[kennung.Etikett](m.GetMetadatei().GetEtiketten())

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

		if err = s.LogWriter.NewOrUpdated(errExists)(sk); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = s.onNewOrUpdatedCommit(sk, false); err != nil {
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

	if ti, err = s.StoreUtil.GetTypenIndex(); err != nil {
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
