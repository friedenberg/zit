package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/alfa/toml"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/golf/persisted_metadatei_format"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
	"github.com/friedenberg/zit/src/kilo/store_util"
)

type CommonStoreBase[
	O objekte.Akte[O],
	OPtr objekte.ObjektePtr[O],
	K schnittstellen.Id[K],
	KPtr schnittstellen.IdPtr[K],
	V any,
	VPtr objekte.VerzeichnissePtr[V, O],
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

	objekte_store.ObjekteSaver
}

type commonStoreBase[
	O objekte.Akte[O],
	OPtr objekte.ObjektePtr[O],
	K schnittstellen.Id[K],
	KPtr schnittstellen.IdPtr[K],
	V any,
	VPtr objekte.VerzeichnissePtr[V, O],
] struct {
	schnittstellen.GattungGetter

	schnittstellen.ObjekteIOFactory

	commonStoreDelegate[O, OPtr, K, KPtr, V, VPtr]

	store_util.StoreUtil

	pool schnittstellen.Pool[
		objekte.Transacted[O, OPtr, K, KPtr, V, VPtr],
		*objekte.Transacted[O, OPtr, K, KPtr, V, VPtr],
	]

	objekte_store.TransactedInflator[O, OPtr, K, KPtr, V, VPtr]

	objekte_store.AkteTextSaver[O, OPtr]

	objekte_store.StoredParseSaver[O, OPtr, K, KPtr]

	objekte_store.TransactedReader[KPtr,
		*objekte.Transacted[O, OPtr, K, KPtr, V, VPtr],
	]

	objekte_store.LogWriter[*objekte.Transacted[O, OPtr, K, KPtr, V, VPtr]]

	persistentMetadateiFormat persisted_metadatei_format.Format

	objekte_store.ObjekteSaver

	textParser metadatei.TextParser
}

func makeCommonStoreBase[
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
	pmf persisted_metadatei_format.Format,
	akteFormat objekte.AkteFormat[O, OPtr],
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
			persisted_metadatei_format.FormatForVersion(
				sa.GetKonfig().GetStoreVersion(),
			),
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
		ObjekteSaver: objekte_store.MakeObjekteSaver(
			of,
			pmf,
		),
		textParser: metadatei.MakeTextParser(
			sa,
			nil, // TODO-P1 make akteFormatter
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
