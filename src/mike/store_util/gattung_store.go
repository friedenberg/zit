package store_util

import (
	"fmt"
	"os"
	"path"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/todo"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/juliett/objekte"
	"github.com/friedenberg/zit/src/lima/objekte_store"
)

type Reindexer interface {
	// updateExternal(objekte.External) error
	ReindexOne(sku.SkuLike) (matcher.Matchable, error)
}

type CommonStoreDelegate[
	O objekte.Akte[O],
	OPtr objekte.AktePtr[O],
	K kennung.KennungLike[K],
	KPtr kennung.KennungLikePtr[K],
] interface {
	AddOne(*sku.Transacted[K, KPtr]) error
	UpdateOne(*sku.Transacted[K, KPtr]) error
}

type CommonStore[
	O objekte.Akte[O],
	OPtr objekte.AktePtr[O],
	K kennung.KennungLike[K],
	KPtr kennung.KennungLikePtr[K],
] struct {
	CommonStoreBase[O, OPtr, K, KPtr]
	AkteFormat objekte.AkteFormat[O, OPtr]
	objekte_store.StoredParseSaver[O, OPtr, K, KPtr]
	objekte_store.CreateOrUpdater[
		OPtr,
		KPtr,
		*sku.Transacted[K, KPtr],
		*objekte.CheckedOut[K, KPtr],
	]
}

func MakeCommonStore[
	O objekte.Akte[O],
	OPtr objekte.AktePtr[O],
	K kennung.KennungLike[K],
	KPtr kennung.KennungLikePtr[K],
](
	gg schnittstellen.GattungGetter,
	delegate CommonStoreDelegate[O, OPtr, K, KPtr],
	sa StoreUtil,
	tr objekte_store.TransactedReader[KPtr,
		*sku.Transacted[K, KPtr]],
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
		StoredParseSaver: objekte_store.MakeStoredParseSaver[O, OPtr, K, KPtr](
			of,
			sa,
			akteFormat,
			sa.GetPersistentMetadateiFormat(),
			objekte_format.Options{IncludeTai: true},
		),
	}

	return
}

func (s *CommonStore[O, OPtr, K, KPtr]) CheckoutOne(
	options CheckoutOptions,
	t *sku.Transacted[K, KPtr],
) (co *objekte.CheckedOut[K, KPtr], err error) {
	todo.Change("add pool")
	co = &objekte.CheckedOut[K, KPtr]{}

	co.Internal = *t
	co.External = t.GetExternal()

	var f *os.File

	p := path.Join(
		s.StoreUtil.GetStandort().Cwd(),
		fmt.Sprintf(
			"%s.%s",
			t.GetKennung(),
			s.StoreUtil.GetKonfig().FileExtensions.GetFileExtensionForGattung(
				t,
			),
		),
	)

	if f, err = files.CreateExclusiveWriteOnly(p); err != nil {
		if errors.IsExist(err) {
			if co.External, err = s.ReadOneExternal(
				sku.ExternalMaybe{
					Kennung: t.GetKennung().KennungPtrClone(),
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

			co.External.Kennung = t.GetKennung()
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	defer errors.DeferredCloser(&err, f)

	if co.External.FDs.Objekte, err = kennung.File(f); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, err = s.AkteFormat.FormatSavedAkte(f, t.GetAkteSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *CommonStore[O, OPtr, K, KPtr]) UpdateManyMetadatei(
	incoming schnittstellen.SetLike[sku.SkuLike],
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
