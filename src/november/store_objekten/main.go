package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/expansion"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/erworben"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/juliett/objekte"
	"github.com/friedenberg/zit/src/kilo/objekte_store"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/mike/store_util"
)

type Store struct {
	store_util.StoreUtil

	protoZettel zettel.ProtoZettel
	konfigStore konfigStore

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
	s.GetKonfig().SetHasChanges(true)

	if !s.GetStandort().GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: "update many metadatei",
		}

		return
	}

	if err = incoming.Each(
		func(mwk *sku.Transacted) (err error) {
			if _, err = s.CreateOrUpdate(
				mwk,
				&mwk.Kennung,
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
	gsWithoutHistory, gsWithHistory := ms.SplitGattungenByHistory()

	wg := iter.MakeErrorWaitGroup()

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

	if err = s.handleNewOrUpdatedCommit(sk, false); err != nil {
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

	var bSha *sha.Sha

	if besty != nil {
		bSha = &besty.Metadatei.Sha
	}

	if err = s.GetEnnui().Add(sk.GetMetadatei(), bSha); err != nil {
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

	f1 := s.ReindexOne

	if err = s.GetBestandsaufnahmeStore().ReadAllSkus(f1); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
