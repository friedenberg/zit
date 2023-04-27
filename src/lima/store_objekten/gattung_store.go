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
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/golf/persisted_metadatei_format"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
	"github.com/friedenberg/zit/src/kilo/store_util"
)

type reindexer interface {
	// updateExternal(objekte.External) error
	ReindexOne(sku.DataIdentity) (kennung.Matchable, error)
}

type CommonStore[
	O objekte.Akte[O],
	OPtr objekte.ObjektePtr[O],
	K schnittstellen.Id[K],
	KPtr schnittstellen.IdPtr[K],
	V any,
	VPtr objekte.VerzeichnissePtr[V, O],
] interface {
	CommonStoreBase[
		O,
		OPtr,
		K,
		KPtr,
		V,
		VPtr,
	]

	objekte_store.CreateOrUpdater[
		OPtr,
		KPtr,
		*objekte.Transacted[
			O,
			OPtr,
			K,
			KPtr,
			V,
			VPtr,
		],
		*objekte.CheckedOut[
			O,
			OPtr,
			K,
			KPtr,
			V,
			VPtr,
		],
	]
}

type commonStoreDelegate[
	O objekte.Akte[O],
	OPtr objekte.ObjektePtr[O],
	K schnittstellen.Id[K],
	KPtr schnittstellen.IdPtr[K],
	V any,
	VPtr objekte.VerzeichnissePtr[V, O],
] interface {
	addOne(*objekte.Transacted[O, OPtr, K, KPtr, V, VPtr]) error
	updateOne(*objekte.Transacted[O, OPtr, K, KPtr, V, VPtr]) error
}

type transacted[T any] interface {
	schnittstellen.Poolable[T]
}

type transactedPtr[T any] interface {
	schnittstellen.PoolablePtr[T]
}

type commonStore[
	O objekte.Akte[O],
	OPtr objekte.ObjektePtr[O],
	K schnittstellen.Id[K],
	KPtr schnittstellen.IdPtr[K],
	V any,
	VPtr objekte.VerzeichnissePtr[V, O],
] struct {
	commonStoreBase[O, OPtr, K, KPtr, V, VPtr]
	AkteFormat objekte.AkteFormat[O, OPtr]
	objekte_store.StoredParseSaver[O, OPtr, K, KPtr]
}

func makeCommonStore[
	O objekte.Akte[O],
	OPtr objekte.ObjektePtr[O],
	K schnittstellen.Id[K],
	KPtr schnittstellen.IdPtr[K],
	V any,
	VPtr objekte.VerzeichnissePtr[V, O],
](
	gg schnittstellen.GattungGetter,
	delegate commonStoreDelegate[O, OPtr, K, KPtr, V, VPtr],
	sa store_util.StoreUtil,
	tr objekte_store.TransactedReader[KPtr,
		*objekte.Transacted[O, OPtr, K, KPtr, V, VPtr]],
	akteFormat objekte.AkteFormat[O, OPtr],
) (s *commonStore[O, OPtr, K, KPtr, V, VPtr], err error) {
	// pool := collections.MakePool[
	// 	objekte.Transacted[O, OPtr, K, KPtr, V, VPtr],
	// 	*objekte.Transacted[O, OPtr, K, KPtr, V, VPtr],
	// ]()

	of := sa.ObjekteReaderWriterFactory(gg)

	csb, err := makeCommonStoreBase[O, OPtr, K, KPtr, V, VPtr](
		gg,
		delegate,
		sa,
		tr,
		persisted_metadatei_format.FormatForVersion(
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
		V,
		VPtr,
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

func (s *commonStore[O, OPtr, K, KPtr, V, VPtr]) CheckoutOne(
	options CheckoutOptions,
	t *objekte.Transacted[O, OPtr, K, KPtr, V, VPtr],
) (co *objekte.CheckedOut[O, OPtr, K, KPtr, V, VPtr], err error) {
	todo.Change("add pool")
	co = &objekte.CheckedOut[O, OPtr, K, KPtr, V, VPtr]{}

	co.Internal = *t
	co.External.Sku = t.Sku.GetExternal()

	var f *os.File

	p := path.Join(
		s.StoreUtil.GetStandort().Cwd(),
		fmt.Sprintf(
			"%s.%s",
			t.Sku.Kennung,
			s.StoreUtil.GetKonfig().FileExtensions.GetFileExtensionForGattung(t),
		),
	)

	if f, err = files.CreateExclusiveWriteOnly(p); err != nil {
		if errors.IsExist(err) {
			if co.External, err = s.ReadOneExternal(
				sku.ExternalMaybe[K, KPtr]{
					Kennung: t.Sku.Kennung,
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

			co.External.Sku.Kennung = t.Sku.Kennung
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
