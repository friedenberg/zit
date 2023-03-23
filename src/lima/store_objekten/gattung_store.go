package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/alfa/toml"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
	"github.com/friedenberg/zit/src/kilo/store_util"
)

type reindexer interface {
	// updateExternal(objekte.External) error
	ReindexOne(sku.DataIdentity) (schnittstellen.Stored, error)
}

type CommonStoreBase[
	O schnittstellen.Objekte[O],
	OPtr schnittstellen.ObjektePtr[O],
	K schnittstellen.Id[K],
	KPtr schnittstellen.IdPtr[K],
	V any,
	VPtr schnittstellen.VerzeichnissePtr[V, O],
] interface {
	reindexer

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
		objekte.External[O, OPtr, K, KPtr],
	]
}

type CommonStore[
	O schnittstellen.Objekte[O],
	OPtr schnittstellen.ObjektePtr[O],
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
	O schnittstellen.Objekte[O],
	OPtr schnittstellen.ObjektePtr[O],
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
	O schnittstellen.Objekte[O],
	OPtr schnittstellen.ObjektePtr[O],
	K schnittstellen.Id[K],
	KPtr schnittstellen.IdPtr[K],
	V any,
	VPtr schnittstellen.VerzeichnissePtr[V, O],
] struct {
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
}

type commonStore[
	O schnittstellen.Objekte[O],
	OPtr schnittstellen.ObjektePtr[O],
	K schnittstellen.Id[K],
	KPtr schnittstellen.IdPtr[K],
	V any,
	VPtr schnittstellen.VerzeichnissePtr[V, O],
] struct {
	commonStoreBase[O, OPtr, K, KPtr, V, VPtr]
	TextFormat schnittstellen.Format[O, OPtr]
	objekte_store.ParseSaver[O, OPtr, K, KPtr]
}

func makeCommonStoreBase[
	O schnittstellen.Objekte[O],
	OPtr schnittstellen.ObjektePtr[O],
	K schnittstellen.Id[K],
	KPtr schnittstellen.IdPtr[K],
	V any,
	VPtr schnittstellen.VerzeichnissePtr[V, O],
](
	delegate commonStoreDelegate[O, OPtr, K, KPtr, V, VPtr],
	sa store_util.StoreUtil,
	tr objekte_store.TransactedReader[KPtr,
		*objekte.Transacted[O, OPtr, K, KPtr, V, VPtr]],
	objekteFormat schnittstellen.Format[O, OPtr],
	textFormat schnittstellen.Format[O, OPtr],
	akteFormatter schnittstellen.Formatter[O, OPtr],
) (s *commonStoreBase[O, OPtr, K, KPtr, V, VPtr], err error) {
	// type T objekte.Transacted[O, OPtr, K, KPtr, V, VPtr]
	// type TPtr *objekte.Transacted[O, OPtr, K, KPtr, V, VPtr]

	pool := collections.MakePool[
		objekte.Transacted[O, OPtr, K, KPtr, V, VPtr],
		*objekte.Transacted[O, OPtr, K, KPtr, V, VPtr],
	]()

	s = &commonStoreBase[O, OPtr, K, KPtr, V, VPtr]{
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
			sa,
			sa,
			objekteFormat,
			textFormat,
			pool,
		),
		AkteTextSaver: objekte_store.MakeAkteTextSaver[
			O,
			OPtr,
		](
			sa,
			akteFormatter,
		),
		TransactedReader: tr,
	}

	return
}

func makeCommonStore[
	O schnittstellen.Objekte[O],
	OPtr schnittstellen.ObjektePtr[O],
	K schnittstellen.Id[K],
	KPtr schnittstellen.IdPtr[K],
	V any,
	VPtr schnittstellen.VerzeichnissePtr[V, O],
](
	delegate commonStoreDelegate[O, OPtr, K, KPtr, V, VPtr],
	sa store_util.StoreUtil,
	tr objekte_store.TransactedReader[KPtr,
		*objekte.Transacted[O, OPtr, K, KPtr, V, VPtr]],
	objekteFormat schnittstellen.Format[O, OPtr],
	textFormat schnittstellen.Format[O, OPtr],
	akteFormatter schnittstellen.Formatter[O, OPtr],
) (s *commonStore[O, OPtr, K, KPtr, V, VPtr], err error) {
	// type T objekte.Transacted[O, OPtr, K, KPtr, V, VPtr]
	// type TPtr *objekte.Transacted[O, OPtr, K, KPtr, V, VPtr]

	pool := collections.MakePool[
		objekte.Transacted[O, OPtr, K, KPtr, V, VPtr],
		*objekte.Transacted[O, OPtr, K, KPtr, V, VPtr],
	]()

	s = &commonStore[
		O,
		OPtr,
		K,
		KPtr,
		V,
		VPtr,
	]{
		commonStoreBase: commonStoreBase[O, OPtr, K, KPtr, V, VPtr]{
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
				sa,
				sa,
				objekteFormat,
				textFormat,
				pool,
			),
			AkteTextSaver: objekte_store.MakeAkteTextSaver[
				O,
				OPtr,
			](
				sa,
				akteFormatter,
			),
			TransactedReader: tr,
		},
		TextFormat: textFormat,
		ParseSaver: objekte_store.MakeParseSaver[
			O,
			OPtr,
			K,
			KPtr,
		](
			sa,
			sa,
			textFormat,
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
	e sku.ExternalMaybe[K, KPtr],
) (t objekte.External[O, OPtr, K, KPtr], err error) {
	if t.Objekte, t.Sku, err = s.ParseAndSaveAkteAndObjekte(
		e,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *commonStoreBase[O, OPtr, K, KPtr, V, VPtr]) ReindexOne(
	sk sku.DataIdentity,
) (o schnittstellen.Stored, err error) {
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
