package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
	"github.com/friedenberg/zit/src/kilo/store_util"
)

type reindexer interface {
	// updateExternal(objekte.External) error
	reindexOne(sku.DataIdentity) (schnittstellen.Stored, error)
}

type CommonStore[
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
}

type transacted[T any] interface {
	schnittstellen.Poolable[T]
}

type transactedPtr[T any] interface {
	schnittstellen.PoolablePtr[T]
}

type commonStore[
	O schnittstellen.Objekte[O],
	OPtr schnittstellen.ObjektePtr[O],
	K schnittstellen.Id[K],
	KPtr schnittstellen.IdPtr[K],
	V any,
	VPtr schnittstellen.VerzeichnissePtr[V, O],
] struct {
	// type T objekte.Transacted[O, OPtr, K, KPtr, V, VPtr]
	// type TPtr *objekte.Transacted[O, OPtr, K, KPtr, V, VPtr]

	store_util.StoreUtil
	pool schnittstellen.Pool[
		objekte.Transacted[O, OPtr, K, KPtr, V, VPtr],
		*objekte.Transacted[O, OPtr, K, KPtr, V, VPtr],
	]

	TextFormat schnittstellen.Format[O, OPtr]

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

	objekte_store.TransactedReader[KPtr,
		*objekte.Transacted[O, OPtr, K, KPtr, V, VPtr],
	]

	objekte_store.LogWriter[*objekte.Transacted[O, OPtr, K, KPtr, V, VPtr]]
}

func makeCommonStore[
	O schnittstellen.Objekte[O],
	OPtr schnittstellen.ObjektePtr[O],
	K schnittstellen.Id[K],
	KPtr schnittstellen.IdPtr[K],
	V any,
	VPtr schnittstellen.VerzeichnissePtr[V, O],
](
	sa store_util.StoreUtil,
	tr objekte_store.TransactedReader[KPtr,
		*objekte.Transacted[O, OPtr, K, KPtr, V, VPtr]],
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
		StoreUtil:  sa,
		pool:       pool,
		TextFormat: textFormat,
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
			nil,
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

func (s *commonStore[O, OPtr, K, KPtr, V, VPtr]) SetLogWriter(
	lw objekte_store.LogWriter[*objekte.Transacted[O, OPtr, K, KPtr, V, VPtr]],
) {
	s.LogWriter = lw
}

func (s *commonStore[O, OPtr, K, KPtr, V, VPtr]) Query(
	m kennung.Matcher,
	f schnittstellen.FuncIter[*objekte.Transacted[O, OPtr, K, KPtr, V, VPtr]],
) (err error) {
	return objekte_store.QueryMethodForMatcher[
		KPtr,
		*objekte.Transacted[O, OPtr, K, KPtr, V, VPtr],
	](s, m, f)
}
