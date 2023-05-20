package store_objekten

import (
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/todo"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/golf/persisted_metadatei_format"
	"github.com/friedenberg/zit/src/golf/transaktion"
	"github.com/friedenberg/zit/src/hotel/erworben"
	"github.com/friedenberg/zit/src/hotel/etikett"
	"github.com/friedenberg/zit/src/hotel/kasten"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
	"github.com/friedenberg/zit/src/hotel/typ"
	"github.com/friedenberg/zit/src/india/bestandsaufnahme"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/kilo/store_util"
)

type Store struct {
	store_util.StoreUtil

	zettelStore  ZettelStore
	typStore     TypStore
	etikettStore EtikettStore
	konfigStore  KonfigStore
	kastenStore  KastenStore

	// Gattungen
	gattungStores     map[schnittstellen.Gattung]gattungStoreLike
	reindexers        map[schnittstellen.Gattung]reindexer
	flushers          map[schnittstellen.Gattung]errors.Flusher
	readers           map[schnittstellen.Gattung]objekte.FuncReaderTransactedLikePtr
	queriers          map[schnittstellen.Gattung]objekte.FuncQuerierTransactedLikePtr
	transactedReaders map[schnittstellen.Gattung]objekte.FuncReaderTransactedLikePtr
	metadateiUpdaters map[schnittstellen.Gattung]objekte_store.UpdaterManyMetadatei

	isReindexing bool
	lock         sync.Locker
}

func Make(
	su store_util.StoreUtil,
	p schnittstellen.Pool[zettel.Transacted, *zettel.Transacted],
) (s *Store, err error) {
	s = &Store{
		lock:      &sync.Mutex{},
		StoreUtil: su,
	}

	su.SetMatchableAdder(s)

	if s.zettelStore, err = makeZettelStore(s.StoreUtil, p); err != nil {
		err = errors.Wrap(err)
		return
	}

	if s.typStore, err = makeTypStore(s.StoreUtil); err != nil {
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

	s.gattungStores = map[schnittstellen.Gattung]gattungStoreLike{
		gattung.Zettel:  s.zettelStore,
		gattung.Typ:     s.typStore,
		gattung.Etikett: s.etikettStore,
		gattung.Konfig:  s.konfigStore,
		gattung.Kasten:  s.kastenStore,
	}

	errors.TodoP1("implement for other gattung")
	s.queriers = map[schnittstellen.Gattung]objekte.FuncQuerierTransactedLikePtr{
		gattung.Zettel: objekte.MakeApplyQueryTransactedLikePtr[*zettel.Transacted](
			s.zettelStore.Query,
		),
		gattung.Typ: objekte.MakeApplyQueryTransactedLikePtr[*typ.Transacted](
			s.typStore.Query,
		),
		gattung.Etikett: objekte.MakeApplyQueryTransactedLikePtr[*etikett.Transacted](
			s.etikettStore.Query,
		),
		gattung.Kasten: objekte.MakeApplyQueryTransactedLikePtr[*kasten.Transacted](
			s.kastenStore.Query,
		),
		gattung.Konfig: objekte.MakeApplyQueryTransactedLikePtr[*erworben.Transacted](
			s.konfigStore.Query,
		),
		// gattung.Bestandsaufnahme: objekte.MakeApplyTransactedLikePtr[*bestandsaufnahme.Objekte](
		// 	s.bestandsaufnahmeStore.ReadAll,
		// ),
	}

	s.readers = map[schnittstellen.Gattung]objekte.FuncReaderTransactedLikePtr{
		gattung.Zettel: objekte.MakeApplyTransactedLikePtr[*zettel.Transacted](
			s.zettelStore.ReadAllSchwanzen,
		),
		gattung.Typ: objekte.MakeApplyTransactedLikePtr[*typ.Transacted](
			s.typStore.ReadAllSchwanzen,
		),
		gattung.Etikett: objekte.MakeApplyTransactedLikePtr[*etikett.Transacted](
			s.etikettStore.ReadAllSchwanzen,
		),
		gattung.Kasten: objekte.MakeApplyTransactedLikePtr[*kasten.Transacted](
			s.kastenStore.ReadAllSchwanzen,
		),
		gattung.Konfig: objekte.MakeApplyTransactedLikePtr[*erworben.Transacted](
			s.konfigStore.ReadAllSchwanzen,
		),
		// gattung.Bestandsaufnahme: objekte.MakeApplyTransactedLikePtr[*bestandsaufnahme.Objekte](
		// 	s.bestandsaufnahmeStore.ReadAll,
		// ),
	}

	s.transactedReaders = map[schnittstellen.Gattung]objekte.FuncReaderTransactedLikePtr{
		gattung.Zettel: objekte.MakeApplyTransactedLikePtr[*zettel.Transacted](
			s.zettelStore.ReadAll,
		),
		gattung.Typ: objekte.MakeApplyTransactedLikePtr[*typ.Transacted](
			s.typStore.ReadAll,
		),
		gattung.Etikett: objekte.MakeApplyTransactedLikePtr[*etikett.Transacted](
			s.etikettStore.ReadAll,
		),
		gattung.Kasten: objekte.MakeApplyTransactedLikePtr[*kasten.Transacted](
			s.kastenStore.ReadAll,
		),
		gattung.Konfig: objekte.MakeApplyTransactedLikePtr[*erworben.Transacted](
			s.konfigStore.ReadAll,
		),
		// gattung.Bestandsaufnahme: objekte.MakeApplyTransactedLikePtr[*bestandsaufnahme.Objekte](
		// 	s.bestandsaufnahmeStore.ReadAll,
		// ),
	}

	s.flushers = make(map[schnittstellen.Gattung]errors.Flusher)

	for g, gs := range s.gattungStores {
		if fl, ok := gs.(errors.Flusher); ok {
			s.flushers[g] = fl
		}
	}

	s.reindexers = make(map[schnittstellen.Gattung]reindexer)

	for g, gs := range s.gattungStores {
		if gs1, ok := gs.(reindexer); ok {
			s.reindexers[g] = gs1
		}
	}

	s.metadateiUpdaters = make(map[schnittstellen.Gattung]objekte_store.UpdaterManyMetadatei)

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
	pmf persisted_metadatei_format.Format,
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
	//			err = errors.Errorf("merge reverts are not yet supported: %s", o)
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
	incoming schnittstellen.Set[metadatei.WithKennung],
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
	ms kennung.MetaSet,
	f schnittstellen.FuncIter[objekte.TransactedLikePtr],
) (err error) {
	if err = ms.All(
		func(g gattung.Gattung, matcher kennung.Matcher) (err error) {
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
	f schnittstellen.FuncIter[objekte.TransactedLikePtr],
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
	f schnittstellen.FuncIter[objekte.TransactedLikePtr],
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

func (s *Store) getReindexFunc() func(sku.DataIdentity) error {
	return func(sk sku.DataIdentity) (err error) {
		var st reindexer
		ok := false

		g := sk.GetGattung()

		if st, ok = s.reindexers[g]; !ok {
			err = gattung.MakeErrUnsupportedGattung(g)
			return
		}

		var o kennung.Matchable

		if o, err = st.ReindexOne(sk); err != nil {
			err = errors.Wrapf(err, "Sku %s", sk)
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
	typenExpanded := kennung.ExpandOneSlice(t, kennung.ExpanderRight)

	for _, t := range typenExpanded {
		if err = s.GetAbbrStore().Typen().Exists(t); err == nil {
			continue
		}

		err = nil

		todo.Change("support inheritance")
		if _, err = s.Typ().CreateOrUpdate(
			typ.MakeObjekte(),
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
	etikettenExpanded := kennung.ExpandOneSlice(e, kennung.ExpanderRight)

	for _, e1 := range etikettenExpanded {
		if err = s.GetAbbrStore().Etiketten().Exists(e1); err == nil {
			continue
		}

		err = nil

		todo.Change("support inheritance")
		if _, err = s.Etikett().CreateOrUpdate(
			etikett.MakeObjekte(),
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
	m kennung.Matchable,
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

	es := collections.SortedValues[kennung.Etikett](m.GetEtiketten())

	for _, e := range es {
		if err = s.addEtikett(e); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s *Store) AddMatchable(m kennung.Matchable) (err error) {
	if err = s.addMatchableTypAndEtikettenIfNecessary(m); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.GetAbbrStore().AddMatchable(m); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

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

	if err = s.StoreUtil.GetKennungIndex().Reset(); err != nil {
		err = errors.Wrapf(err, "failed to reset index kennung")
		return
	}

	f1 := s.getReindexFunc()

	// if s.StoreUtil.GetKonfig().UseBestandsaufnahme {
	// } else {
	f := func(t *transaktion.Transaktion) (err error) {
		errors.Out().Printf("%s/%s: %s", t.Time.Kopf(), t.Time.Schwanz(), t.Time)

		if err = t.Skus.Each(
			func(sk sku.SkuLike) (err error) {
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

	f2 := func(t *bestandsaufnahme.Transacted) (err error) {
		if err = t.Akte.Skus.Each(
			func(sk sku.Sku) (err error) {
				return f1(sk)
			},
		); err != nil {
			err = errors.Wrapf(
				err,
				"Bestandsaufnahme: %s",
				t.GetKennung(),
			)

			return
		}

		return
	}

	if err = s.GetBestandsaufnahmeStore().ReadAll(f2); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
