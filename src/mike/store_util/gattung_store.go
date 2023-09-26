package store_util

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/kilo/objekte_store"
)

type Reindexer interface {
	// updateExternal(objekte.External) error
	// ReindexOne(*sku.Transacted) (matcher.Matchable, error)
}

type CommonStore struct {
	CommonStoreBase
	cou objekte_store.CreateOrUpdater
}

func MakeCommonStore(
	gg schnittstellen.GattungGetter,
	sa StoreUtil,
	tr objekte_store.TransactedReader,
	cou objekte_store.CreateOrUpdater,
) (s *CommonStore, err error) {
	csb, err := MakeCommonStoreBase(
		gg,
		sa,
		tr,
		objekte_format.FormatForVersion(
			sa.GetStandort().GetKonfig().GetStoreVersion(),
		),
	)
	if err != nil {
		err = errors.Wrap(err)
		return
	}

	s = &CommonStore{
		CommonStoreBase: *csb,
		cou:             cou,
	}

	return
}

func (s *CommonStore) UpdateManyMetadatei(
	incoming sku.TransactedSet,
) (err error) {
	if !s.GetStandort().GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: "update many metadatei",
		}

		return
	}

	if err = incoming.EachPtr(
		func(mwk *sku.Transacted) (err error) {
			if !mwk.GetGattung().EqualsGattung(s) {
				return
			}

			if _, err = s.cou.CreateOrUpdate(
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
