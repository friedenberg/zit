package store_objekten

import (
	"fmt"
	"os"
	"path"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/alfa/toml"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/todo"
	"github.com/friedenberg/zit/src/charlie/collections"
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

type CommonStoreBase[
	O objekte.Objekte[O],
	OPtr objekte.ObjektePtr[O],
	K schnittstellen.Id[K],
	KPtr schnittstellen.IdPtr[K],
	V any,
	VPtr schnittstellen.VerzeichnissePtr[V, O],
] interface {
	reindexer

	CheckoutOne(
		CheckoutOptions,
		*objekte.Transacted[O, OPtr, K, KPtr, V, VPtr],
	) (*objekte.CheckedOut[O, OPtr, K, KPtr, V, VPtr], error)

	objekte_store.TransactedLogger[*objekte.Transacted[
		O,
		OPtr,
		K,
		KPtr,
		V,
		VPtr,
	]]

	objekte_store.Querier[
		KPtr,
		*objekte.Transacted[
			O,
			OPtr,
			K,
			KPtr,
			V,
			VPtr,
		],
	]

	objekte_store.AkteTextSaver[
		O,
		OPtr,
	]

	objekte_store.TransactedInflator[
		O,
		OPtr,
		K,
		KPtr,
		V,
		VPtr,
	]

	objekte_store.Inheritor[*objekte.Transacted[
		O,
		OPtr,
		K,
		KPtr,
		V,
		VPtr,
	]]

	objekte_store.ExternalReader[
		sku.ExternalMaybe[K, KPtr],
		*objekte.Transacted[O, OPtr, K, KPtr, V, VPtr],
		objekte.External[O, OPtr, K, KPtr],
	]

	schnittstellen.ObjekteIOFactory
}

type CommonStore[
	O objekte.Objekte[O],
	OPtr objekte.ObjektePtr[O],
	K schnittstellen.Id[K],
	KPtr schnittstellen.IdPtr[K],
	V any,
	VPtr schnittstellen.VerzeichnissePtr[V, O],
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
	O objekte.Objekte[O],
	OPtr objekte.ObjektePtr[O],
	K schnittstellen.Id[K],
	KPtr schnittstellen.IdPtr[K],
	V any,
	VPtr schnittstellen.VerzeichnissePtr[V, O],
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

type commonStoreBase[
	O objekte.Objekte[O],
	OPtr objekte.ObjektePtr[O],
	K schnittstellen.Id[K],
	KPtr schnittstellen.IdPtr[K],
	V any,
	VPtr schnittstellen.VerzeichnissePtr[V, O],
] struct {
	schnittstellen.GattungGetter

	schnittstellen.ObjekteIOFactory

	commonStoreDelegate[O, OPtr, K, KPtr, V, VPtr]

	store_util.StoreUtil

	pool schnittstellen.Pool[
		objekte.Transacted[O, OPtr, K, KPtr, V, VPtr],
		*objekte.Transacted[O, OPtr, K, KPtr, V, VPtr],
	]

	objekte_store.TransactedInflator[
		O,
		OPtr,
		K,
		KPtr,
		V,
		VPtr,
	]

	objekte_store.AkteTextSaver[
		O,
		OPtr,
	]

	objekte_store.ParseSaver[
		O,
		OPtr,
		K,
		KPtr,
	]

	objekte_store.TransactedReader[KPtr,
		*objekte.Transacted[O, OPtr, K, KPtr, V, VPtr],
	]

	objekte_store.LogWriter[*objekte.Transacted[O, OPtr, K, KPtr, V, VPtr]]

	persistentMetadateiFormat persisted_metadatei_format.Format
}

type commonStore[
	O objekte.Objekte[O],
	OPtr objekte.ObjektePtr[O],
	K schnittstellen.Id[K],
	KPtr schnittstellen.IdPtr[K],
	V any,
	VPtr schnittstellen.VerzeichnissePtr[V, O],
] struct {
	commonStoreBase[O, OPtr, K, KPtr, V, VPtr]
	AkteFormat objekte_store.AkteFormat[O, OPtr]
	objekte_store.ParseSaver[O, OPtr, K, KPtr]
}

func makeCommonStoreBase[
	O objekte.Objekte[O],
	OPtr objekte.ObjektePtr[O],
	K schnittstellen.Id[K],
	KPtr schnittstellen.IdPtr[K],
	V any,
	VPtr schnittstellen.VerzeichnissePtr[V, O],
](
	gg schnittstellen.GattungGetter,
	delegate commonStoreDelegate[O, OPtr, K, KPtr, V, VPtr],
	sa store_util.StoreUtil,
	tr objekte_store.TransactedReader[KPtr,
		*objekte.Transacted[O, OPtr, K, KPtr, V, VPtr]],
	pmf persisted_metadatei_format.Format,
	akteFormat objekte_store.AkteFormat[O, OPtr],
) (s *commonStoreBase[O, OPtr, K, KPtr, V, VPtr], err error) {
	// type T objekte.Transacted[O, OPtr, K, KPtr, V, VPtr]
	// type TPtr *objekte.Transacted[O, OPtr, K, KPtr, V, VPtr]

	pool := collections.MakePool[
		objekte.Transacted[O, OPtr, K, KPtr, V, VPtr],
		*objekte.Transacted[O, OPtr, K, KPtr, V, VPtr],
	]()

	of := sa.ObjekteReaderWriterFactory(gg)

	s = &commonStoreBase[O, OPtr, K, KPtr, V, VPtr]{
		GattungGetter:       gg,
		ObjekteIOFactory:    of,
		commonStoreDelegate: delegate,
		StoreUtil:           sa,
		pool:                pool,
		TransactedInflator: objekte_store.MakeTransactedInflator[
			O,
			OPtr,
			K,
			KPtr,
			V,
			VPtr,
		](
			of,
			sa,
			persisted_metadatei_format.V0{},
			akteFormat,
			pool,
		),
		AkteTextSaver: objekte_store.MakeAkteTextSaver[
			O,
			OPtr,
		](
			sa,
			akteFormat,
		),
		TransactedReader:          tr,
		persistentMetadateiFormat: pmf,
	}

	return
}

func makeCommonStore[
	O objekte.Objekte[O],
	OPtr objekte.ObjektePtr[O],
	K schnittstellen.Id[K],
	KPtr schnittstellen.IdPtr[K],
	V any,
	VPtr schnittstellen.VerzeichnissePtr[V, O],
](
	gg schnittstellen.GattungGetter,
	delegate commonStoreDelegate[O, OPtr, K, KPtr, V, VPtr],
	sa store_util.StoreUtil,
	tr objekte_store.TransactedReader[KPtr,
		*objekte.Transacted[O, OPtr, K, KPtr, V, VPtr]],
	akteFormat objekte_store.AkteFormat[O, OPtr],
) (s *commonStore[O, OPtr, K, KPtr, V, VPtr], err error) {
	pool := collections.MakePool[
		objekte.Transacted[O, OPtr, K, KPtr, V, VPtr],
		*objekte.Transacted[O, OPtr, K, KPtr, V, VPtr],
	]()

	of := sa.ObjekteReaderWriterFactory(gg)

	s = &commonStore[
		O,
		OPtr,
		K,
		KPtr,
		V,
		VPtr,
	]{
		commonStoreBase: commonStoreBase[O, OPtr, K, KPtr, V, VPtr]{
			GattungGetter:       gg,
			ObjekteIOFactory:    of,
			commonStoreDelegate: delegate,
			StoreUtil:           sa,
			pool:                pool,
			TransactedInflator: objekte_store.MakeTransactedInflator[
				O,
				OPtr,
				K,
				KPtr,
				V,
				VPtr,
			](
				of,
				sa,
				persisted_metadatei_format.V0{},
				akteFormat,
				pool,
			),
			AkteTextSaver: objekte_store.MakeAkteTextSaver[
				O,
				OPtr,
			](
				sa,
				akteFormat,
			),
			TransactedReader: tr,
		},
		AkteFormat: akteFormat,
		ParseSaver: objekte_store.MakeParseSaver[
			O,
			OPtr,
			K,
			KPtr,
		](
			of,
			sa,
			akteFormat,
			sa.GetPersistentMetadateiFormat(),
		),
	}

	return
}

func (s *commonStoreBase[O, OPtr, K, KPtr, V, VPtr]) SetLogWriter(
	lw objekte_store.LogWriter[*objekte.Transacted[O, OPtr, K, KPtr, V, VPtr]],
) {
	s.LogWriter = lw
}

func (s *commonStoreBase[O, OPtr, K, KPtr, V, VPtr]) Query(
	m kennung.Matcher,
	f schnittstellen.FuncIter[*objekte.Transacted[O, OPtr, K, KPtr, V, VPtr]],
) (err error) {
	return objekte_store.QueryMethodForMatcher[
		KPtr,
		*objekte.Transacted[O, OPtr, K, KPtr, V, VPtr],
	](s, m, f)
}

func (s *commonStore[O, OPtr, K, KPtr, V, VPtr]) ReadOneExternal(
	em sku.ExternalMaybe[K, KPtr],
	t *objekte.Transacted[O, OPtr, K, KPtr, V, VPtr],
) (e objekte.External[O, OPtr, K, KPtr], err error) {
	// support akte
	todo.Implement()
	if e.Objekte, e.Sku, err = s.ParseAndSaveAkteAndObjekte(
		em,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *commonStoreBase[O, OPtr, K, KPtr, V, VPtr]) ReindexOne(
	sk sku.DataIdentity,
) (o kennung.Matchable, err error) {
	var t *objekte.Transacted[O, OPtr, K, KPtr, V, VPtr]

	if t, err = s.InflateFromDataIdentity(sk); err != nil {
		if errors.Is(err, toml.Error{}) {
			err = nil
			return
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	o = t

	if t.IsNew() {
		s.LogWriter.New(t)
		if err = s.addOne(t); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		s.LogWriter.Updated(t)
		if err = s.updateOne(t); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s *commonStoreBase[O, OPtr, K, KPtr, V, VPtr]) Inherit(
	t *objekte.Transacted[O, OPtr, K, KPtr, V, VPtr],
) (err error) {
	if t == nil {
		err = errors.Errorf("trying to inherit nil %T", t)
		return
	}

	errors.Log().Printf("inheriting %s", t.Sku.ObjekteSha)

	s.StoreUtil.CommitTransacted(t)

	old, _ := s.ReadOne(&t.Sku.Kennung)

	if old == nil || old.Less(*t) {
		s.addOne(t)
	}

	if t.IsNew() {
		s.LogWriter.New(t)
	} else {
		s.LogWriter.Updated(t)
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

	if _, err = s.AkteFormat.Format(f, &t.Objekte); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
