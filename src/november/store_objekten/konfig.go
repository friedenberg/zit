package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/erworben"
	"github.com/friedenberg/zit/src/juliett/objekte"
	"github.com/friedenberg/zit/src/kilo/objekte_store"
	"github.com/friedenberg/zit/src/mike/store_util"
)

type konfigStore struct {
	store_util.StoreUtil

	akteFormat objekte.AkteFormat[erworben.Akte, *erworben.Akte]
	objekte_store.LogWriter
}

func (s *konfigStore) GetAkteFormat() objekte.AkteFormat[erworben.Akte, *erworben.Akte] {
	return s.akteFormat
}

func (s konfigStore) Update(
	sh schnittstellen.ShaLike,
) (kt *sku.Transacted, err error) {
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

	kt = sku.GetTransactedPool().Get()

	if err = kt.Kennung.SetWithKennung(&kennung.Konfig{}); err != nil {
		err = errors.Wrap(err)
		return
	}

	kt.SetTai(s.GetTai())
	kt.SetAkteSha(sh)

	// TODO-P3 refactor into reusable
	if mutter != nil {
		if err = kt.Metadatei.Verzeichnisse.Mutter.SetShaLike(
			&mutter.GetMetadatei().Verzeichnisse.Sha,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	err = sku.CalculateAndSetSha(kt, s.GetPersistentMetadateiFormat(),
		objekte_format.Options{IncludeTai: true},
	)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if mutter != nil && kt.Metadatei.EqualsSansTai(&mutter.Metadatei) {
		if err = kt.SetFromSkuLike(mutter); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = s.Unchanged(kt); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

  if err = iter.Chain(
    kt,
    s.AddVerzeichnisse,
    s.CommitUpdatedTransacted,
    s.AddMatchable,
    s.Updated,
  ); err != nil {
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

	tt1 := &s.StoreUtil.GetKonfig().Sku

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
			s.GetPersistentMetadateiFormat(),
			objekte_format.Options{IncludeTai: true},
		)

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
