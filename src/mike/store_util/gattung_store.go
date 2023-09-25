package store_util

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/kilo/objekte_store"
)

type Reindexer interface {
	// updateExternal(objekte.External) error
	ReindexOne(*sku.Transacted) (matcher.Matchable, error)
}

type CommonStoreDelegate interface {
	AddOne(*sku.Transacted) error
	UpdateOne(*sku.Transacted) error
}

type CommonStore struct {
	CommonStoreBase
	objekte_store.CreateOrUpdater
}

func MakeCommonStore(
	gg schnittstellen.GattungGetter,
	delegate CommonStoreDelegate,
	sa StoreUtil,
	tr objekte_store.TransactedReader,
) (s *CommonStore, err error) {
	if delegate == nil {
		panic("delegate was nil")
	}

	csb, err := MakeCommonStoreBase(
		gg,
		delegate,
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

			if _, err = s.CreateOrUpdater.CreateOrUpdate(
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
