package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/erworben"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/juliett/objekte"
	"github.com/friedenberg/zit/src/kilo/objekte_store"
	"github.com/friedenberg/zit/src/mike/store_util"
)

type konfigStore struct {
	*store_util.CommonStore[
		erworben.Akte,
		*erworben.Akte,
	]

	akteFormat objekte.AkteFormat[erworben.Akte, *erworben.Akte]
}

func (s *konfigStore) GetAkteFormat() objekte.AkteFormat[erworben.Akte, *erworben.Akte] {
	return s.akteFormat
}

func makeKonfigStore(
	sa store_util.StoreUtil,
) (s *konfigStore, err error) {
	s = &konfigStore{
		akteFormat: objekte_store.MakeAkteFormat[erworben.Akte, *erworben.Akte](
			objekte.MakeTextParserIgnoreTomlErrors[erworben.Akte](
				sa.GetStandort(),
			),
			objekte.ParsedAkteTomlFormatter[erworben.Akte]{},
			sa.GetStandort(),
		),
	}

	s.CommonStore, err = store_util.MakeCommonStore[
		erworben.Akte,
		*erworben.Akte,
	](
		gattung.Konfig,
		s,
		sa,
		s,
	)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s konfigStore) Flush() (err error) {
	return
}

func (s konfigStore) AddOne(t *erworben.Transacted) (err error) {
	s.StoreUtil.GetKonfigPtr().SetTransacted(t, s.GetAkten().GetKonfigV0())
	return
}

func (s konfigStore) UpdateOne(t *erworben.Transacted) (err error) {
	s.StoreUtil.GetKonfigPtr().SetTransacted(t, s.GetAkten().GetKonfigV0())
	return
}

func (s konfigStore) Update(
	sh schnittstellen.ShaLike,
) (kt *erworben.Transacted, err error) {
	if !s.StoreUtil.GetStandort().GetLockSmith().IsAcquired() {
		err = errors.Wrap(
			objekte_store.ErrLockRequired{Operation: "update konfig"},
		)
		return
	}

	var mutter *sku.Transacted

	if mutter, err = s.ReadOne(&kennung.Konfig{}); err != nil {
		if errors.Is(err, objekte_store.ErrNotFound{}) {
			mutter = nil
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	kt = &erworben.Transacted{}

	kt.Kennung = kennung.Kennung2{KennungPtr: &kennung.Konfig{}}

	kt.SetTai(s.StoreUtil.GetTai())
	kt.SetAkteSha(sh)

	// TODO-P3 refactor into reusable
	if mutter != nil {
		kt.Metadatei.Verzeichnisse.Mutter = mutter.GetMetadatei().Verzeichnisse.Sha
	}

	err = sku.CalculateAndSetSha(kt, s.StoreUtil.GetPersistentMetadateiFormat(),
		objekte_format.Options{IncludeTai: true},
	)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if mutter != nil && kt.Metadatei.EqualsSansTai(mutter.GetMetadatei()) {
		if err = kt.SetFromSkuLike(mutter); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = s.LogWriter.Unchanged(kt); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	s.StoreUtil.CommitUpdatedTransacted(kt)

	if err = s.StoreUtil.GetKonfigPtr().SetTransacted(kt, s.GetAkten().GetKonfigV0()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.StoreUtil.AddMatchable(kt); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.LogWriter.Updated(kt); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *konfigStore) ReadAllSchwanzen(
	w schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	var k *sku.Transacted

	if k, err = i.ReadOne(&kennung.Konfig{}); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = w(k); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *konfigStore) ReadAll(
	w schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	eachSku := func(sk *sku.Transacted) (err error) {
		if sk.GetGattung() != gattung.Konfig {
			return
		}

		if err = w(sk); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	if err = s.StoreUtil.GetBestandsaufnahmeStore().ReadAllSkus(eachSku); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s konfigStore) ReadOne(
	k schnittstellen.StringerGattungGetter,
) (tt *sku.Transacted, err error) {
	var k1 kennung.Konfig

	if err = k1.Set(k.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	tt1 := s.StoreUtil.GetKonfig().Sku

	if tt1.GetTai().IsEmpty() {
		err = errors.Wrap(objekte_store.ErrNotFound{Id: k1})
		return
	}

	tt = sku.GetTransactedPool().Get()

	if err = tt.SetFromSkuLike(tt1); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !tt.GetTai().IsEmpty() {
		err = sku.CalculateAndSetSha(
			tt,
			s.StoreUtil.GetPersistentMetadateiFormat(),
			objekte_format.Options{IncludeTai: true},
		)

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s *konfigStore) ReindexOne(
	te *sku.Transacted,
) (o matcher.Matchable, err error) {
	o = te

	if err = s.UpdateOne(te); err != nil {
		err = errors.Wrap(err)
		return
	}

	s.LogWriter.Updated(te)

	return
}
