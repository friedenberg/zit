package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/golf/transaktion"
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
	gattungStores     map[schnittstellen.Gattung]GattungStore
	reindexers        map[schnittstellen.Gattung]reindexer
	flushers          map[schnittstellen.Gattung]errors.Flusher
	readers           map[schnittstellen.Gattung]objekte.FuncReaderTransactedLike
	queriers          map[schnittstellen.Gattung]objekte.FuncQuerierTransactedLike
	transactedReaders map[schnittstellen.Gattung]objekte.FuncReaderTransactedLike
}

func Make(
	su store_util.StoreUtil,
	p schnittstellen.Pool[zettel.Transacted, *zettel.Transacted],
) (s *Store, err error) {
	s = &Store{
		StoreUtil: su,
	}

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

	s.gattungStores = map[schnittstellen.Gattung]GattungStore{
		gattung.Zettel:  s.zettelStore,
		gattung.Typ:     s.typStore,
		gattung.Etikett: s.etikettStore,
		gattung.Konfig:  s.konfigStore,
		gattung.Kasten:  s.kastenStore,
	}

	errors.TodoP0("implement for other gattung")
	s.queriers = map[schnittstellen.Gattung]objekte.FuncQuerierTransactedLike{
		gattung.Zettel: objekte.MakeApplyQueryTransactedLike[*zettel.Transacted](
			s.zettelStore.Query,
		),
		gattung.Typ: objekte.MakeApplyQueryTransactedLike[*typ.Transacted](
			s.typStore.Query,
		),
		// gattung.Typ: objekte.MakeApplyTransactedLike[*typ.Transacted](
		// 	s.typStore.ReadAllSchwanzen,
		// ),
		// gattung.Etikett: objekte.MakeApplyTransactedLike[*etikett.Transacted](
		// 	s.etikettStore.ReadAllSchwanzen,
		// ),
		// gattung.Kasten: objekte.MakeApplyTransactedLike[*kasten.Transacted](
		// 	s.kastenStore.ReadAllSchwanzen,
		// ),
		// gattung.Konfig:           objekte.MakeApplyTransactedLike[*konfig.Transacted](
		// s.konfigStore.ReadAllSchwanzen,
		// ),
		// gattung.Bestandsaufnahme: objekte.MakeApplyTransactedLike[*bestandsaufnahme.Objekte](
		// 	s.bestandsaufnahmeStore.ReadAll,
		// ),
	}

	s.readers = map[schnittstellen.Gattung]objekte.FuncReaderTransactedLike{
		gattung.Zettel: objekte.MakeApplyTransactedLike[*zettel.Transacted](
			s.zettelStore.ReadAllSchwanzen,
		),
		gattung.Typ: objekte.MakeApplyTransactedLike[*typ.Transacted](
			s.typStore.ReadAllSchwanzen,
		),
		gattung.Etikett: objekte.MakeApplyTransactedLike[*etikett.Transacted](
			s.etikettStore.ReadAllSchwanzen,
		),
		gattung.Kasten: objekte.MakeApplyTransactedLike[*kasten.Transacted](
			s.kastenStore.ReadAllSchwanzen,
		),
		// gattung.Konfig:           objekte.MakeApplyTransactedLike[*konfig.Transacted](
		// s.konfigStore.ReadAllSchwanzen,
		// ),
		// gattung.Bestandsaufnahme: objekte.MakeApplyTransactedLike[*bestandsaufnahme.Objekte](
		// 	s.bestandsaufnahmeStore.ReadAll,
		// ),
	}

	s.transactedReaders = map[schnittstellen.Gattung]objekte.FuncReaderTransactedLike{
		gattung.Zettel: objekte.MakeApplyTransactedLike[*zettel.Transacted](
			s.zettelStore.ReadAll,
		),
		gattung.Typ: objekte.MakeApplyTransactedLike[*typ.Transacted](
			s.typStore.ReadAll,
		),
		gattung.Etikett: objekte.MakeApplyTransactedLike[*etikett.Transacted](
			s.etikettStore.ReadAll,
		),
		gattung.Kasten: objekte.MakeApplyTransactedLike[*kasten.Transacted](
			s.kastenStore.ReadAll,
		),
		// gattung.Konfig:           objekte.MakeApplyTransactedLike[*konfig.Transacted](
		// s.konfigStore.ReadAllSchwanzen,
		// ),
		// gattung.Bestandsaufnahme: objekte.MakeApplyTransactedLike[*bestandsaufnahme.Objekte](
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
	errors.TodoP0("implement for Bestandsaufnahme")

	if !s.StoreUtil.GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: "revert",
		}

		return
	}

	// tzs = zettel.MakeMutableSetUnique(t.Skus.Len())

	// t.Skus.Each(
	//	func(o sku.SkuLike) (err error) {
	//		var h *hinweis.Hinweis
	//		ok := false

	//		if h, ok = o.GetId().(*hinweis.Hinweis); !ok {
	//			//TODO
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

	return
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
	if _, err = s.GetBestandsaufnahmeStore().Create(
		s.StoreUtil.GetBestandsaufnahme(),
	); err != nil {
		if errors.Is(err, bestandsaufnahme.ErrEmpty) {
			errors.Log().Printf("Bestandsaufnahme was empty")
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}
	errors.Log().Printf("done saving Bestandsaufnahme")

	if err = s.StoreUtil.GetTransaktionStore().WriteTransaktion(); err != nil {
		err = errors.Wrapf(err, "failed to write transaction")
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

func (s *Store) Query(
	ms kennung.MetaSet,
	f schnittstellen.FuncIter[objekte.TransactedLike],
) (err error) {
	if err = ms.All(
		func(g gattung.Gattung, ids kennung.Set) (err error) {
			r, ok := s.queriers[g]

			if !ok {
				return
			}

			if err = r(ids, f); err != nil {
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
	f schnittstellen.FuncIter[objekte.TransactedLike],
) (err error) {
	chErr := make(chan error, gs.Len())

	for g, s1 := range s.readers {
		if !gs.ContainsKey(g.GetGattungString()) {
			continue
		}

		go func(s1 objekte.FuncReaderTransactedLike) {
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
	f schnittstellen.FuncIter[objekte.TransactedLike],
) (err error) {
	chErr := make(chan error, gs.Len())

	for g, s1 := range s.transactedReaders {
		if !gs.ContainsKey(g.GetGattungString()) {
			continue
		}

		go func(s1 objekte.FuncReaderTransactedLike) {
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

		var o schnittstellen.Stored

		if o, err = st.reindexOne(sk); err != nil {
			err = errors.Wrapf(err, "Sku %s", sk)
			return
		}

		if err = s.GetAbbrStore().AddStoredAbbreviation(o); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}

func (s *Store) Reindex() (err error) {
	if !s.GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: "reindex",
		}

		return
	}

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

	f2 := func(t *bestandsaufnahme.Objekte) (err error) {
		if err = t.Akte.Skus.Each(
			func(sk sku.Sku2) (err error) {
				return f1(sk)
			},
		); err != nil {
			err = errors.Wrapf(
				err,
				"Bestandsaufnahme: %s",
				t.Tai,
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
