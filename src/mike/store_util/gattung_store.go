package store_util

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/juliett/objekte"
	"github.com/friedenberg/zit/src/lima/objekte_store"
)

type Reindexer interface {
	// updateExternal(objekte.External) error
	ReindexOne(*sku.Transacted) (matcher.Matchable, error)
}

type CommonStoreDelegate interface {
	AddOne(*sku.Transacted) error
	UpdateOne(*sku.Transacted) error
}

type CommonStore[
	O objekte.Akte[O],
	OPtr objekte.AktePtr[O],
	K kennung.KennungLike[K],
	KPtr kennung.KennungLikePtr[K],
] struct {
	CommonStoreBase[O, OPtr]
	AkteFormat objekte.AkteFormat[O, OPtr]
	objekte_store.StoredParseSaver[O, OPtr]
	objekte_store.CreateOrUpdater
}

func MakeCommonStore[
	O objekte.Akte[O],
	OPtr objekte.AktePtr[O],
	K kennung.KennungLike[K],
	KPtr kennung.KennungLikePtr[K],
](
	gg schnittstellen.GattungGetter,
	delegate CommonStoreDelegate,
	sa StoreUtil,
	tr objekte_store.TransactedReader,
	akteFormat objekte.AkteFormat[O, OPtr],
) (s *CommonStore[O, OPtr, K, KPtr], err error) {
	// pool := collections.MakePool[
	// 	objekte.Transacted[O, OPtr, K, KPtr, ],
	// 	*objekte.Transacted[O, OPtr, K, KPtr, ],
	// ]()

	if delegate == nil {
		panic("delegate was nil")
	}

	of := sa.ObjekteReaderWriterFactory(gg)

	csb, err := MakeCommonStoreBase[O, OPtr, K, KPtr](
		gg,
		delegate,
		sa,
		tr,
		objekte_format.FormatForVersion(
			sa.GetStoreVersion(),
		),
		akteFormat,
	)
	if err != nil {
		err = errors.Wrap(err)
		return
	}

	s = &CommonStore[
		O,
		OPtr,
		K,
		KPtr,
	]{
		CommonStoreBase: *csb,
		AkteFormat:      akteFormat,
		StoredParseSaver: objekte_store.MakeStoredParseSaver[O, OPtr](
			of,
			sa,
			akteFormat,
			sa.GetPersistentMetadateiFormat(),
			objekte_format.Options{IncludeTai: true},
		),
	}

	return
}

func (s *CommonStore[O, OPtr, K, KPtr]) UpdateManyMetadatei(
	incoming sku.TransactedSet,
) (err error) {
	if !s.StoreUtil.GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: "update many metadatei",
		}

		return
	}

	if err = incoming.EachPtr(
		func(mwk *sku.Transacted) (err error) {
			k := mwk.GetKennungLike()

			var ke K
			kep := KPtr(&ke)

			switch kt := k.(type) {
			case K:
				kep = &kt

			case KPtr:
				kep = kt

			case kennung.Kennung2, *kennung.Kennung2:
				if !kt.GetGattung().EqualsGattung(ke.GetGattung()) {
					return
				}

				if err = kep.Set(kt.String()); err != nil {
					err = errors.Wrap(err)
					return
				}

			default:
				return
			}

			if _, err = s.CreateOrUpdater.CreateOrUpdate(
				mwk,
				kep,
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
