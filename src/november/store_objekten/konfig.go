package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/kilo/objekte_store"
)

func (s *Store) UpdateKonfig(
	sh schnittstellen.ShaLike,
) (kt *sku.Transacted, err error) {
	return s.CreateOrUpdateAkte(
		nil,
		&kennung.Konfig{},
		sh,
	)

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
		if err = kt.Metadatei.Mutter.SetShaLike(
			&mutter.GetMetadatei().Sha,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	err = sku.CalculateAndSetSha(kt, s.GetPersistentMetadateiFormat(),
		objekte_format.Options{Tai: true},
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

	if err = s.handleUpdated(
		kt,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
