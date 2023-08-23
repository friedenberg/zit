package store_objekten

import (
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/bravo/todo"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/golf/kennung_index"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/erworben"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/india/transacted"
	"github.com/friedenberg/zit/src/india/transaktion"
	"github.com/friedenberg/zit/src/juliett/objekte"
	"github.com/friedenberg/zit/src/kilo/bestandsaufnahme"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/lima/objekte_store"
	"github.com/friedenberg/zit/src/mike/store_util"
)

type Store struct {
	store_util.StoreUtil

	zettelStore  ZettelStore
	typStore     TypStore
	etikettStore EtikettStore
	konfigStore  KonfigStore
	kastenStore  KastenStore

	// Gattungen
	gattungStores     map[schnittstellen.GattungLike]gattungStoreLike
	reindexers        map[schnittstellen.GattungLike]reindexer
	flushers          map[schnittstellen.GattungLike]errors.Flusher
	readers           map[schnittstellen.GattungLike]objekte.FuncReaderTransactedLikePtr
	queriers          map[schnittstellen.GattungLike]objekte.FuncQuerierTransactedLikePtr
	transactedReaders map[schnittstellen.GattungLike]objekte.FuncReaderTransactedLikePtr
	metadateiUpdaters map[schnittstellen.GattungLike]objekte_store.UpdaterManyMetadatei

	isReindexing bool
	lock         sync.Locker
}

func Make(
	su store_util.StoreUtil,
	p schnittstellen.Pool[transacted.Zettel, *transacted.Zettel],
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

	s.gattungStores = map[schnittstellen.GattungLike]gattungStoreLike{
		gattung.Zettel:  s.zettelStore,
		gattung.Typ:     s.typStore,
		gattung.Etikett: s.etikettStore,
		gattung.Konfig:  s.konfigStore,
		gattung.Kasten:  s.kastenStore,
	}

	errors.TodoP1("implement for other gattung")
	s.queriers = map[schnittstellen.GattungLike]objekte.FuncQuerierTransactedLikePtr{
		gattung.Zettel: objekte.MakeApplyQueryTransactedLikePtr[*transacted.Zettel](
			s.zettelStore.Query,
		),
		gattung.Typ: objekte.MakeApplyQueryTransactedLikePtr[*transacted.Typ](
			s.typStore.Query,
		),
		gattung.Etikett: objekte.MakeApplyQueryTransactedLikePtr[*transacted.Etikett](
			s.etikettStore.Query,
		),
		gattung.Kasten: objekte.MakeApplyQueryTransactedLikePtr[*transacted.Kasten](
			s.kastenStore.Query,
		),
		gattung.Konfig: objekte.MakeApplyQueryTransactedLikePtr[*erworben.Transacted](
			s.konfigStore.Query,
		),
		// gattung.Bestandsaufnahme:
		// objekte.MakeApplyTransactedLikePtr[*bestandsaufnahme.Objekte](
		// 	s.bestandsaufnahmeStore.ReadAll,
		// ),
	}

	s.readers = map[schnittstellen.GattungLike]objekte.FuncReaderTransactedLikePtr{
		gattung.Zettel: objekte.MakeApplyTransactedLikePtr[*transacted.Zettel](
			s.zettelStore.ReadAllSchwanzen,
		),
		gattung.Typ: objekte.MakeApplyTransactedLikePtr[*transacted.Typ](
			s.typStore.ReadAllSchwanzen,
		),
		gattung.Etikett: objekte.MakeApplyTransactedLikePtr[*transacted.Etikett](
			s.etikettStore.ReadAllSchwanzen,
		),
		gattung.Kasten: objekte.MakeApplyTransactedLikePtr[*transacted.Kasten](
			s.kastenStore.ReadAllSchwanzen,
		),
		gattung.Konfig: objekte.MakeApplyTransactedLikePtr[*erworben.Transacted](
			s.konfigStore.ReadAllSchwanzen,
		),
		// gattung.Bestandsaufnahme:
		// objekte.MakeApplyTransactedLikePtr[*bestandsaufnahme.Objekte](
		// 	s.bestandsaufnahmeStore.ReadAll,
		// ),
	}

	s.transactedReaders = map[schnittstellen.GattungLike]objekte.FuncReaderTransactedLikePtr{
		gattung.Zettel: objekte.MakeApplyTransactedLikePtr[*transacted.Zettel](
			s.zettelStore.ReadAll,
		),
		gattung.Typ: objekte.MakeApplyTransactedLikePtr[*transacted.Typ](
			s.typStore.ReadAll,
		),
		gattung.Etikett: objekte.MakeApplyTransactedLikePtr[*transacted.Etikett](
			s.etikettStore.ReadAll,
		),
		gattung.Kasten: objekte.MakeApplyTransactedLikePtr[*transacted.Kasten](
			s.kastenStore.ReadAll,
		),
		gattung.Konfig: objekte.MakeApplyTransactedLikePtr[*erworben.Transacted](
			s.konfigStore.ReadAll,
		),
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

	s.reindexers = make(map[schnittstellen.GattungLike]reindexer)

	for g, gs := range s.gattungStores {
		if gs1, ok := gs.(reindexer); ok {
			s.reindexers[g] = gs1
		}
	}

	s.metadateiUpdaters = make(
		map[schnittstellen.GattungLike]objekte_store.UpdaterManyMetadatei,
	)

	for g, gs := range s.gattungStores {
		if gs1, ok := gs.(objekte_store.UpdaterManyMetadatei); ok {
			s.metadateiUpdaters[g] = gs1
		}
	}

	return
}

func (s *Store) GetGattungInheritors(
	ofg schnittstellen.ObjekteReaderFactoryGetter,
	af schnittstellen.AkteReaderFactory,
	pmf objekte_format.Format,
) (out map[gattung.Gattung]objekte_store.TransactedInheritor) {
	out = make(map[gattung.Gattung]objekte_store.TransactedInheritor)

	for g1, gs := range s.gattungStores {
		g := gattung.Make(g1)
		of := ofg.ObjekteReaderFactory(g)
		out[g] = gs.GetInheritor(of, af, pmf)
	}

	return
}

func (s *Store) Zettel() ZettelStore {
	return s.zettelStore
}

func (s *Store) Typ() TypStore {
	return s.typStore
}

func (s *Store) Etikett() EtikettStore {
	return s.etikettStore
}

func (s *Store) Konfig() KonfigStore {
	return s.konfigStore
}

func (s *Store) Kasten() KastenStore {
	return s.kastenStore
}

func (s Store) RevertTransaktion(
	t *transaktion.Transaktion,
) (tzs zettel.MutableSet, err error) {
	err = errors.Implement()
	return

	// if !s.StoreUtil.GetLockSmith().IsAcquired() {
	// 	err = objekte_store.ErrLockRequired{
	// 		Operation: "revert",
	// 	}

	// 	return
	// }

	// tzs = zettel.MakeMutableSetUnique(t.Skus.Len())

	// t.Skus.Each(
	//	func(o sku.SkuLike) (err error) {
	//		var h *hinweis.Hinweis
	//		ok := false

	//		if h, ok = o.GetId().(*hinweis.Hinweis); !ok {
	//			//TODO-P4
	//			return
	//		}

	//		if !o.GetMutter()[1].IsZero() {
	// 			err = errors.Errorf("merge reverts are not yet supported: %s", o)
	//			return
	//		}

	//		errors.Log().Print(o)

	//		var chain []*zettel.Transacted

	//		if chain, err = s.zettelStore.AllInChain(*h); err != nil {
	//			err = errors.Wrap(err)
	//			return
	//		}

	//		var tz *zettel.Transacted

	//		for _, someTz := range chain {
	//			errors.Log().Print(someTz)
	//			if someTz.Sku.Schwanz == o.GetMutter()[0] {
	//				tz = someTz
	//				break
	//			}
	//		}

	//		if tz.Sku.ObjekteSha.IsNull() {
	//			err = errors.Errorf("zettel not found in index!: %#v", o)
	//			return
	//		}

	//		if tz, err = s.zettelStore.Update(
	//			&tz.Objekte,
	//			&tz.Sku.Kennung,
	//		); err != nil {
	//			err = errors.Wrap(err)
	//			return
	//		}

	//		tzs.Add(tz)

	//		return
	//	},
	//)

	// return
}

func (s Store) Flush() (err error) {
	if !s.GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: "flush",
		}

		return
	}

	if s.GetKonfig().DryRun {
		return
	}

	errors.Log().Printf("saving Bestandsaufnahme")
	ba := s.StoreUtil.GetBestandsaufnahmeAkte()
	if err = s.GetBestandsaufnahmeStore().Create(&ba); err != nil {
		if errors.Is(err, bestandsaufnahme.ErrEmpty) {
			errors.Log().Printf("Bestandsaufnahme was empty")
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}
	errors.Log().Printf("done saving Bestandsaufnahme")

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
	incoming schnittstellen.SetLike[sku.SkuLike],
) (err error) {
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
	ms matcher.MetaSet,
	f schnittstellen.FuncIter[sku.SkuLikePtr],
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

func (s *Store) ReadAllSchwanzen(
	gs gattungen.Set,
	f schnittstellen.FuncIter[sku.SkuLikePtr],
) (err error) {
	chErr := make(chan error, gs.Len())

	for g, s1 := range s.readers {
		if !gs.ContainsKey(g.GetGattungString()) {
			continue
		}

		go func(s1 objekte.FuncReaderTransactedLikePtr) {
			var subErr error

			defer func() {
				chErr <- subErr
			}()

			subErr = s1(f)
		}(s1)
	}

	for i := 0; i < gs.Len(); i++ {
		err = errors.MakeMulti(err, <-chErr)
	}

	return
}

func (s *Store) ReadAll(
	gs gattungen.Set,
	f schnittstellen.FuncIter[sku.SkuLikePtr],
) (err error) {
	chErr := make(chan error, gs.Len())

	for g, s1 := range s.transactedReaders {
		if !gs.ContainsKey(g.GetGattungString()) {
			continue
		}

		go func(s1 objekte.FuncReaderTransactedLikePtr) {
			var subErr error

			defer func() {
				chErr <- subErr
			}()

			subErr = s1(f)
		}(s1)
	}

	for i := 0; i < gs.Len(); i++ {
		err = errors.MakeMulti(err, <-chErr)
	}

	return
}

func (s *Store) GetReindexFunc(
	ti kennung_index.KennungIndex[kennung.Typ, *kennung.Typ],
) func(sku.SkuLikePtr) error {
	return func(sk sku.SkuLikePtr) (err error) {
		var st reindexer
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
		if _, err = s.Typ().CreateOrUpdate(
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
	if !s.GetLockSmith().IsAcquired() {
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

	// if s.StoreUtil.GetKonfig().UseBestandsaufnahme {
	// } else {
	f := func(t *transaktion.Transaktion) (err error) {
		errors.Out().Printf(
			"%s/%s: %s",
			t.Time.Kopf(),
			t.Time.Schwanz(),
			t.Time,
		)

		if err = t.Skus.Each(
			func(sk sku.SkuLikePtr) (err error) {
				return f1(sk)
			},
		); err != nil {
			err = errors.Wrapf(
				err,
				"Transaktion: %s/%s: %s",
				t.Time.Kopf(),
				t.Time.Schwanz(),
				t.Time,
			)

			return
		}

		return
	}

	if err = s.GetTransaktionStore().ReadAllTransaktions(f); err != nil {
		if errors.IsNotExist(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	// }

	if err = s.GetBestandsaufnahmeStore().ReadAllSkus(f1); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
