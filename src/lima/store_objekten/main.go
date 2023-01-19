package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/age"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/standort"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/echo/hinweis"
	"github.com/friedenberg/zit/src/echo/ts"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/golf/transaktion"
	"github.com/friedenberg/zit/src/hotel/bestandsaufnahme"
	"github.com/friedenberg/zit/src/hotel/etikett"
	"github.com/friedenberg/zit/src/hotel/typ"
	"github.com/friedenberg/zit/src/india/konfig"
	"github.com/friedenberg/zit/src/juliett/zettel"
)

type Store struct {
	common

	zettelStore  *zettelStore
	typStore     TypStore
	etikettStore EtikettStore
	konfigStore  KonfigStore

	//Gattungen
	gattungStores     map[schnittstellen.Gattung]GattungStore
	reindexers        map[schnittstellen.Gattung]reindexer
	flushers          map[schnittstellen.Gattung]errors.Flusher
	readers           map[schnittstellen.Gattung]objekte.FuncReaderTransactedLike
	transactedReaders map[schnittstellen.Gattung]objekte.FuncReaderTransactedLike
}

func Make(
	lockSmith LockSmith,
	a age.Age,
	k *konfig.Compiled,
	st standort.Standort,
	p *collections.Pool[zettel.Transacted],
) (s *Store, err error) {
	s = &Store{
		common: common{
			LockSmith: lockSmith,
			Age:       a,
			konfig:    k,
			Standort:  st,
		},
	}

	t := ts.Now()
	ta := ts.NowTai()

	for {
		p := s.TransaktionPath(t)

		if !files.Exists(p) {
			break
		}

		t.MoveForwardIota()
	}

	s.common.Transaktion = transaktion.MakeTransaktion(t)
	s.common.Bestandsaufnahme = &bestandsaufnahme.Objekte{
		Tai: ta,
		Akte: bestandsaufnahme.Akte{
			Skus: sku.MakeSku2Heap(),
		},
	}

	if s.common.Abbr, err = newIndexAbbr(
		&s.common,
		st.DirVerzeichnisse("Abbr"),
	); err != nil {
		err = errors.Wrapf(err, "failed to init abbr index")
		return
	}

	if s.zettelStore, err = makeZettelStore(&s.common, p); err != nil {
		err = errors.Wrap(err)
		return
	}

	if s.typStore, err = makeTypStore(&s.common); err != nil {
		err = errors.Wrap(err)
		return
	}

	if s.etikettStore, err = makeEtikettStore(&s.common); err != nil {
		err = errors.Wrap(err)
		return
	}

	if s.konfigStore, err = makeKonfigStore(&s.common); err != nil {
		err = errors.Wrap(err)
		return
	}

	if s.bestandsaufnahmeStore, err = bestandsaufnahme.MakeStore(
		s.common.Standort,
		&s.common,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	s.gattungStores = map[schnittstellen.Gattung]GattungStore{
		gattung.Zettel:  s.zettelStore,
		gattung.Typ:     s.typStore,
		gattung.Etikett: s.etikettStore,
		gattung.Konfig:  s.konfigStore,
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

func (s *Store) Abbr() *indexAbbr {
	return s.common.Abbr
}

func (s *Store) Zettel() *zettelStore {
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

func (s *Store) CurrentTransaktionTime() ts.Time {
	return s.common.Transaktion.Time
}

func (s Store) RevertTransaktion(
	t *transaktion.Transaktion,
) (tzs zettel.MutableSet, err error) {
	errors.TodoP0("implement")

	if !s.common.LockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "revert",
		}

		return
	}

	tzs = zettel.MakeMutableSetUnique(t.Skus.Len())

	t.Skus.Each(
		func(o sku.SkuLike) (err error) {
			var h *hinweis.Hinweis
			ok := false

			if h, ok = o.GetId().(*hinweis.Hinweis); !ok {
				//TODO
				return
			}

			if !o.GetMutter()[1].IsZero() {
				err = errors.Errorf("merge reverts are not yet supported: %s", o)
				return
			}

			errors.Log().Print(o)

			var chain []*zettel.Transacted

			if chain, err = s.zettelStore.AllInChain(*h); err != nil {
				err = errors.Wrap(err)
				return
			}

			var tz *zettel.Transacted

			for _, someTz := range chain {
				errors.Log().Print(someTz)
				if someTz.Sku.Schwanz == o.GetMutter()[0] {
					tz = someTz
					break
				}
			}

			if tz.Sku.ObjekteSha.IsNull() {
				err = errors.Errorf("zettel not found in index!: %#v", o)
				return
			}

			if tz, err = s.zettelStore.Update(
				&tz.Objekte,
				&tz.Sku.Kennung,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			tzs.Add(tz)

			return
		},
	)

	return
}

func (s Store) Flush() (err error) {
	if !s.common.LockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "flush",
		}

		return
	}

	if s.common.Konfig().DryRun {
		return
	}

	if _, err = s.bestandsaufnahmeStore.Create(s.common.Bestandsaufnahme); err != nil {
		if errors.Is(err, bestandsaufnahme.ErrEmpty) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	//TODO-P2 add Bestandsaufnahme to Transaktion

	if err = s.writeTransaktion(); err != nil {
		err = errors.Wrapf(err, "failed to write transaction")
		return
	}

	for _, fl := range s.flushers {
		if err = fl.Flush(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = s.common.Abbr.Flush(); err != nil {
		errors.Err().Print(err)
		err = errors.Wrapf(err, "failed to flush abbr index")
		return
	}

	return
}

func (s *Store) ReadAllSchwanzen(
	gs gattungen.Set,
	f collections.WriterFunc[objekte.TransactedLike],
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
	f collections.WriterFunc[objekte.TransactedLike],
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
			err = errors.Wrapf(gattung.ErrUnsupportedGattung, "Gattung: %s", g)
			return
		}

		var o schnittstellen.Stored

		if o, err = st.reindexOne(sk); err != nil {
			err = errors.Wrapf(err, "Sku %s", sk)
			return
		}

		if err = s.common.Abbr.addStoredAbbreviation(o); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}

func (s *Store) Reindex() (err error) {
	if !s.common.LockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "reindex",
		}

		return
	}

	if err = s.common.Standort.ResetVerzeichnisse(); err != nil {
		err = errors.Wrap(err)
		return
	}

	//TODO-P3 move to zettelStore
	if err = s.zettelStore.indexKennung.reset(); err != nil {
		err = errors.Wrapf(err, "failed to reset index kennung")
		return
	}

	f1 := s.getReindexFunc()

	if s.common.Konfig().UseBestandsaufnahme {
		f := func(t *bestandsaufnahme.Objekte) (err error) {
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

		if err = s.bestandsaufnahmeStore.ReadAll(f); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
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

		if err = s.ReadAllTransaktions(f); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	return
}
