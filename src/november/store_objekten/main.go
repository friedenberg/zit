package store_objekten

import (
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/bravo/todo"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/golf/kennung_index"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/kilo/objekte_store"
	"github.com/friedenberg/zit/src/mike/store_util"
)

type CreateOrUpdator interface {
	CreateOrUpdateCheckedOut(
		co *sku.CheckedOut,
	) (transactedPtr *sku.Transacted, err error)

	CreateOrUpdateAkte(
		mg metadatei.Getter,
		kennungPtr *kennung.Kennung2,
		sh schnittstellen.ShaLike,
	) (transactedPtr *sku.Transacted, err error)
}

type Store struct {
	store_util.StoreUtil

	zettelStore  *zettelStore
	typStore     *typStore
	etikettStore *etikettStore
	konfigStore  *konfigStore
	kastenStore  *kastenStore

	CreateOrUpdator CreateOrUpdator

	objekte_store.LogWriter

	// Gattungen
	gattungStores     map[schnittstellen.GattungLike]store_util.GattungStoreLike
	reindexers        map[schnittstellen.GattungLike]store_util.Reindexer
	flushers          map[schnittstellen.GattungLike]errors.Flusher
	readers           map[schnittstellen.GattungLike]matcher.FuncReaderTransactedLikePtr
	queriers          map[schnittstellen.GattungLike]matcher.FuncSigilTransactedLikePtr
	transactedReaders map[schnittstellen.GattungLike]matcher.FuncReaderTransactedLikePtr
	metadateiUpdaters map[schnittstellen.GattungLike]objekte_store.UpdaterManyMetadatei

	isReindexing bool
	lock         sync.Locker
}

func Make(
	su store_util.StoreUtil,
	p schnittstellen.Pool[sku.Transacted, *sku.Transacted],
) (s *Store, err error) {
	s = &Store{
		lock:      &sync.Mutex{},
		StoreUtil: su,
	}

	su.SetMatchableAdder(s)

	if s.typStore, err = makeTypStore(s.StoreUtil); err != nil {
		err = errors.Wrap(err)
		return
	}

	if s.zettelStore, err = makeZettelStore(s.StoreUtil, p, s.typStore); err != nil {
		err = errors.Wrap(err)
		return
	}

	if s.etikettStore, err = makeEtikettStore(s.StoreUtil); err != nil {
		err = errors.Wrap(err)
		return
	}

	if s.konfigStore, err = makeKonfigStore(s.StoreUtil); err != nil {
		err = errors.Wrap(err)
		return
	}

	if s.kastenStore, err = makeKastenStore(s.StoreUtil); err != nil {
		err = errors.Wrap(err)
		return
	}

	s.gattungStores = map[schnittstellen.GattungLike]store_util.GattungStoreLike{
		gattung.Zettel:  s.zettelStore,
		gattung.Typ:     s.typStore,
		gattung.Etikett: s.etikettStore,
		gattung.Konfig:  s.konfigStore,
		gattung.Kasten:  s.kastenStore,
	}

	errors.TodoP1("implement for other gattung")
	s.queriers = map[schnittstellen.GattungLike]matcher.FuncSigilTransactedLikePtr{
		gattung.Zettel:  s.zettelStore.Query,
		gattung.Typ:     s.typStore.Query,
		gattung.Etikett: s.etikettStore.Query,
		gattung.Kasten:  s.kastenStore.Query,
		gattung.Konfig:  s.konfigStore.Query,
		// gattung.Bestandsaufnahme:
		// objekte.MakeApplyTransactedLikePtr[*bestandsaufnahme.Objekte](
		// 	s.bestandsaufnahmeStore.ReadAll,
		// ),
	}

	s.readers = map[schnittstellen.GattungLike]matcher.FuncReaderTransactedLikePtr{
		gattung.Zettel:  s.zettelStore.ReadAllSchwanzen,
		gattung.Typ:     s.typStore.ReadAllSchwanzen,
		gattung.Etikett: s.etikettStore.ReadAllSchwanzen,
		gattung.Kasten:  s.kastenStore.ReadAllSchwanzen,
		gattung.Konfig:  s.konfigStore.ReadAllSchwanzen,
		// gattung.Bestandsaufnahme:
		// objekte.MakeApplyTransactedLikePtr[*bestandsaufnahme.Objekte](
		// 	s.bestandsaufnahmeStore.ReadAll,
		// ),
	}

	s.transactedReaders = map[schnittstellen.GattungLike]matcher.FuncReaderTransactedLikePtr{
		gattung.Zettel:  s.zettelStore.ReadAll,
		gattung.Typ:     s.typStore.ReadAll,
		gattung.Etikett: s.etikettStore.ReadAll,
		gattung.Kasten:  s.kastenStore.ReadAll,
		gattung.Konfig:  s.konfigStore.ReadAll,
		// gattung.Bestandsaufnahme:
		// objekte.MakeApplyTransactedLikePtr[*bestandsaufnahme.Objekte](
		// 	s.bestandsaufnahmeStore.ReadAll,
		// ),
	}

	s.flushers = make(map[schnittstellen.GattungLike]errors.Flusher)

	for g, gs := range s.gattungStores {
		if fl, ok := gs.(errors.Flusher); ok {
			s.flushers[g] = fl
		}
	}

	s.reindexers = make(map[schnittstellen.GattungLike]store_util.Reindexer)

	for g, gs := range s.gattungStores {
		if gs1, ok := gs.(store_util.Reindexer); ok {
			s.reindexers[g] = gs1
		}
	}

	s.metadateiUpdaters = map[schnittstellen.GattungLike]objekte_store.UpdaterManyMetadatei{
		gattung.Zettel:  s.zettelStore,
		gattung.Typ:     s.typStore,
		gattung.Etikett: s.etikettStore,
		gattung.Kasten:  s.kastenStore,
		// gattung.Konfig:  s.konfigStore,
	}

	s.CreateOrUpdator = objekte_store.MakeCreateOrUpdate2(
		s,
		s.GetStandort().GetLockSmith(),
		s.GetStandort(),
		s,
		objekte_store.CreateOrUpdateDelegate{
			New:       s.onNew,
			Updated:   s.onUpdated,
			Unchanged: s.onUnchanged,
		},
		s,
		s.GetPersistentMetadateiFormat(),
		objekte_format.Options{IncludeTai: true},
		s.StoreUtil,
		sku.GetTransactedPool(),
	)

	return
}

func (s *Store) SetLogWriter(
	lw objekte_store.LogWriter,
) {
	s.LogWriter = lw
	s.zettelStore.SetLogWriter(lw)
	s.konfigStore.SetLogWriter(lw)
	s.typStore.SetLogWriter(lw)
	s.etikettStore.SetLogWriter(lw)
	s.kastenStore.SetLogWriter(lw)
}

func (s *Store) Zettel() *zettelStore {
	return s.zettelStore
}

func (s *Store) Typ() *typStore {
	return s.typStore
}

func (s *Store) Etikett() *etikettStore {
	return s.etikettStore
}

func (s *Store) Konfig() *konfigStore {
	return s.konfigStore
}

func (s *Store) Kasten() *kastenStore {
	return s.kastenStore
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
	s.GetKonfigPtr().SetHasChanges(true)
	todo.Optimize() // parallelize
	for _, umm := range s.metadateiUpdaters {
		if err = umm.UpdateManyMetadatei(incoming); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s *Store) Query(
	ms matcher.Query,
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	if err = ms.All(
		func(g gattung.Gattung, matcher matcher.MatcherSigil) (err error) {
			r, ok := s.queriers[g]

			if !ok {
				return
			}

			if err = r(matcher, f); err != nil {
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

func (s *Store) GetReindexFunc(
	ti kennung_index.KennungIndex[kennung.Typ, *kennung.Typ],
) func(*sku.Transacted) error {
	return func(sk *sku.Transacted) (err error) {
		var st store_util.Reindexer
		ok := false

		g := sk.GetGattung()

		if st, ok = s.reindexers[g]; !ok {
			err = gattung.MakeErrUnsupportedGattung(g)
			return
		}

		var o matcher.Matchable

		if o, err = st.ReindexOne(sk); err != nil {
			err = errors.Wrapf(err, "Sku %s", sk)
			return
		}

		if err = ti.StoreOne(o.GetTyp()); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = s.GetAbbrStore().AddMatchable(o); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}

func (s *Store) addTyp(
	t kennung.Typ,
) (err error) {
	typenExpanded := kennung.ExpandOneSlice(&t, kennung.ExpanderRight)

	for _, t := range typenExpanded {
		if err = s.GetAbbrStore().Typen().Exists(t); err == nil {
			continue
		}

		err = nil

		todo.Change("support inheritance")
		if _, err = s.typStore.CreateOrUpdate(
			nil,
			&t,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s *Store) addEtikett(
	e kennung.Etikett,
) (err error) {
	etikettenExpanded := kennung.ExpandOneSlice(&e, kennung.ExpanderRight)

	for _, e1 := range etikettenExpanded {
		if err = s.GetAbbrStore().Etiketten().Exists(e1); err == nil {
			continue
		}

		err = nil

		todo.Change("support inheritance")
		if _, err = s.Etikett().CreateOrUpdate(
			nil,
			&e1,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s *Store) addMatchableTypAndEtikettenIfNecessary(
	m matcher.Matchable,
) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	// // TODO-P2 support other true gattung
	// if !gattung.Zettel.EqualsAny(m.GetGattung()) {
	// 	return
	// }

	t := m.GetTyp()

	if err = s.addTyp(t); err != nil {
		err = errors.Wrap(err)
		return
	}

	es := iter.SortedValues[kennung.Etikett](m.GetMetadatei().GetEtiketten())

	for _, e := range es {
		if err = s.addEtikett(e); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s *Store) AddMatchable(m matcher.Matchable) (err error) {
	if err = s.addMatchableTypAndEtikettenIfNecessary(m); err != nil {
		err = errors.Wrap(err)
		return
	}

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

	if err = s.StoreUtil.GetStandort().ResetVerzeichnisse(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var ti kennung_index.KennungIndex[kennung.Typ, *kennung.Typ]

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

	f1 := s.GetReindexFunc(ti)

	if err = s.GetBestandsaufnahmeStore().ReadAllSkus(f1); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
