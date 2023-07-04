package store_objekten

import (
	"fmt"
	"os"
	"path"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/todo"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/hotel/objekte"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
	"github.com/friedenberg/zit/src/kilo/store_util"
)

type reindexer interface {
	// updateExternal(objekte.External) error
	ReindexOne(sku.SkuLike) (kennung.Matchable, error)
}

type CommonStore[
	O objekte.Akte[O],
	OPtr objekte.AktePtr[O],
	K kennung.KennungLike[K],
	KPtr kennung.KennungLikePtr[K],
] interface {
	CommonStoreBase[O, OPtr, K, KPtr]

	objekte_store.CreateOrUpdater[
		OPtr,
		KPtr,
		*objekte.Transacted[O, OPtr, K, KPtr],
		*objekte.CheckedOut[O, OPtr, K, KPtr],
	]
}

type commonStoreDelegate[
	O objekte.Akte[O],
	OPtr objekte.AktePtr[O],
	K kennung.KennungLike[K],
	KPtr kennung.KennungLikePtr[K],
] interface {
	addOne(*objekte.Transacted[O, OPtr, K, KPtr]) error
	updateOne(*objekte.Transacted[O, OPtr, K, KPtr]) error
}

type transacted[T any] interface {
	schnittstellen.Poolable[T]
}

type transactedPtr[T any] interface {
	schnittstellen.PoolablePtr[T]
}

type commonStore[
	O objekte.Akte[O],
	OPtr objekte.AktePtr[O],
	K kennung.KennungLike[K],
	KPtr kennung.KennungLikePtr[K],
] struct {
	commonStoreBase[O, OPtr, K, KPtr]
	AkteFormat objekte.AkteFormat[O, OPtr]
	objekte_store.StoredParseSaver[O, OPtr, K, KPtr]
	objekte_store.CreateOrUpdater[
		OPtr,
		KPtr,
		*objekte.Transacted[O, OPtr, K, KPtr],
		*objekte.CheckedOut[O, OPtr, K, KPtr],
	]
}

func makeCommonStore[
	O objekte.Akte[O],
	OPtr objekte.AktePtr[O],
	K kennung.KennungLike[K],
	KPtr kennung.KennungLikePtr[K],
](
	gg schnittstellen.GattungGetter,
	delegate commonStoreDelegate[O, OPtr, K, KPtr],
	sa store_util.StoreUtil,
	tr objekte_store.TransactedReader[KPtr,
		*objekte.Transacted[O, OPtr, K, KPtr]],
	akteFormat objekte.AkteFormat[O, OPtr],
) (s *commonStore[O, OPtr, K, KPtr], err error) {
	// pool := collections.MakePool[
	// 	objekte.Transacted[O, OPtr, K, KPtr, ],
	// 	*objekte.Transacted[O, OPtr, K, KPtr, ],
	// ]()

	if delegate == nil {
		panic("delegate was nil")
	}

	of := sa.ObjekteReaderWriterFactory(gg)

	csb, err := makeCommonStoreBase[O, OPtr, K, KPtr](
		gg,
		delegate,
		sa,
		tr,
		objekte_format.FormatForVersion(
			sa.GetKonfig().GetStoreVersion(),
		),
		akteFormat,
	)
	if err != nil {
		err = errors.Wrap(err)
		return
	}

	s = &commonStore[
		O,
		OPtr,
		K,
		KPtr,
	]{
		commonStoreBase: *csb,
		AkteFormat:      akteFormat,
		StoredParseSaver: objekte_store.MakeStoredParseSaver[O, OPtr, K, KPtr](
			of,
			sa,
			akteFormat,
			sa.GetPersistentMetadateiFormat(),
		),
	}

	return
}

func (s *commonStore[O, OPtr, K, KPtr]) CheckoutOne(
	options CheckoutOptions,
	t *objekte.Transacted[O, OPtr, K, KPtr],
) (co *objekte.CheckedOut[O, OPtr, K, KPtr], err error) {
	todo.Change("add pool")
	co = &objekte.CheckedOut[O, OPtr, K, KPtr]{}

	co.Internal = *t
	co.External.Sku = t.Sku.GetExternal()

	var f *os.File

	p := path.Join(
		s.StoreUtil.GetStandort().Cwd(),
		fmt.Sprintf(
			"%s.%s",
			t.Sku.GetKennung(),
			s.StoreUtil.GetKonfig().FileExtensions.GetFileExtensionForGattung(
				t,
			),
		),
	)

	if f, err = files.CreateExclusiveWriteOnly(p); err != nil {
		if errors.IsExist(err) {
			if co.External, err = s.ReadOneExternal(
				sku.ExternalMaybe[K, KPtr]{
					Kennung: t.Sku.GetKennung(),
					FDs: sku.ExternalFDs{
						Objekte: kennung.FD{
							Path: p,
						},
					},
				},
				t,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			co.External.Sku.Kennung = t.Sku.GetKennung()
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	defer errors.DeferredCloser(&err, f)

	if co.External.Sku.FDs.Objekte, err = kennung.File(f); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, err = s.AkteFormat.FormatSavedAkte(f, t.GetAkteSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *commonStore[O, OPtr, K, KPtr]) UpdateManyMetadatei(
	incoming schnittstellen.Set[sku.SkuLike],
) (err error) {
	if !s.StoreUtil.GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: "update many metadatei",
		}

		return
	}

	if err = incoming.Each(
		func(mwk sku.SkuLike) (err error) {
			var ke K
			ok := false

			if ke, ok = mwk.GetKennungLike().(K); !ok {
				return
			}

			kep := KPtr(&ke)

			var old *objekte.Transacted[O, OPtr, K, KPtr]

			if old, err = s.ReadOne(kep); err != nil {
				err = errors.Wrap(err)
				return
			}

			if _, err = s.CreateOrUpdater.CreateOrUpdate(
				&old.Akte,
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
