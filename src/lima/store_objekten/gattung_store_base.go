package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/alfa/toml"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/hotel/objekte"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
	"github.com/friedenberg/zit/src/kilo/store_util"
)

type gattungStoreLike interface {
	reindexer
	schnittstellen.ObjekteIOFactory
	objekte_store.ObjekteSaver
	GetInheritor(
		schnittstellen.ObjekteReaderFactory,
		schnittstellen.AkteReaderFactory,
		objekte_format.Format,
	) objekte_store.TransactedInheritor
}

type CommonStoreBase[
	O objekte.Akte[O],
	OPtr objekte.AktePtr[O],
	K kennung.KennungLike[K],
	KPtr kennung.KennungLikePtr[K],
] interface {
	gattungStoreLike

	CheckoutOne(
		CheckoutOptions,
		*objekte.Transacted[O, OPtr, K, KPtr],
	) (*objekte.CheckedOut[O, OPtr, K, KPtr], error)

	objekte_store.TransactedLogger[*objekte.Transacted[
		O,
		OPtr,
		K,
		KPtr,
	]]

	objekte_store.Querier[
		KPtr,
		*objekte.Transacted[
			O,
			OPtr,
			K,
			KPtr,
		],
	]

	objekte_store.AkteTextSaver[O, OPtr]

	objekte_store.TransactedInflator[O, OPtr, K, KPtr]

	objekte_store.Inheritor[*objekte.Transacted[
		O,
		OPtr,
		K,
		KPtr,
	]]

	objekte_store.ExternalReader[
		*sku.ExternalMaybe[K, KPtr],
		*objekte.Transacted[O, OPtr, K, KPtr],
		objekte.External[O, OPtr, K, KPtr],
	]
}

type commonStoreBase[
	O objekte.Akte[O],
	OPtr objekte.AktePtr[O],
	K kennung.KennungLike[K],
	KPtr kennung.KennungLikePtr[K],
] struct {
	schnittstellen.GattungGetter

	schnittstellen.ObjekteIOFactory

	delegate commonStoreDelegate[O, OPtr, K, KPtr]

	store_util.StoreUtil

	pool schnittstellen.Pool[
		objekte.Transacted[O, OPtr, K, KPtr],
		*objekte.Transacted[O, OPtr, K, KPtr],
	]

	objekte_store.TransactedInflator[O, OPtr, K, KPtr]

	objekte_store.AkteTextSaver[O, OPtr]

	objekte_store.StoredParseSaver[O, OPtr, K, KPtr]

	objekte_store.TransactedReader[KPtr,
		*objekte.Transacted[O, OPtr, K, KPtr],
	]

	objekte_store.LogWriter[objekte.TransactedLikePtr]

	persistentMetadateiFormat objekte_format.Format

	objekte_store.ObjekteSaver

	textParser metadatei.TextParser
}

func makeCommonStoreBase[
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
	pmf objekte_format.Format,
	akteFormat objekte.AkteFormat[O, OPtr],
) (s *commonStoreBase[O, OPtr, K, KPtr], err error) {
	// type T objekte.Transacted[O, OPtr, K, KPtr, ]
	// type TPtr *objekte.Transacted[O, OPtr, K, KPtr, ]

	if delegate == nil {
		panic("delegate was nil")
	}

	pool := collections.MakePool[
		objekte.Transacted[O, OPtr, K, KPtr],
		*objekte.Transacted[O, OPtr, K, KPtr],
	]()

	of := sa.ObjekteReaderWriterFactory(gg)

	s = &commonStoreBase[O, OPtr, K, KPtr]{
		GattungGetter:    gg,
		ObjekteIOFactory: of,
		delegate:         delegate,
		StoreUtil:        sa,
		pool:             pool,
		TransactedInflator: objekte_store.MakeTransactedInflator[
			O,
			OPtr,
			K,
			KPtr,
		](
			sa.GetStoreVersion(),
			of,
			sa,
			objekte_format.FormatForVersion(
				sa.GetStoreVersion(),
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

func (s *commonStoreBase[O, OPtr, K, KPtr]) SetLogWriter(
	lw objekte_store.LogWriter[objekte.TransactedLikePtr],
) {
	s.LogWriter = lw
}

func (s *commonStoreBase[O, OPtr, K, KPtr]) Query(
	m kennung.MatcherSigil,
	f schnittstellen.FuncIter[*objekte.Transacted[O, OPtr, K, KPtr]],
) (err error) {
	return objekte_store.QueryMethodForMatcher[
		KPtr,
		*objekte.Transacted[O, OPtr, K, KPtr],
	](s, m, f)
}

func (s *commonStoreBase[O, OPtr, K, KPtr]) ReindexOne(
	sk sku.SkuLike,
) (o kennung.Matchable, err error) {
	var t *objekte.Transacted[O, OPtr, K, KPtr]

	if t, err = s.InflateFromSku(sk); err != nil {
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
		if err = s.delegate.addOne(t); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		s.LogWriter.Updated(t)
		if err = s.delegate.updateOne(t); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s *commonStoreBase[O, OPtr, K, KPtr]) Inherit(
	t *objekte.Transacted[O, OPtr, K, KPtr],
) (err error) {
	if t == nil {
		err = errors.Errorf("trying to inherit nil %T", t)
		return
	}

	errors.Log().Printf("inheriting %s", t.Sku.ObjekteSha)

	s.StoreUtil.AddMatchable(t)
	s.StoreUtil.CommitTransacted(t)

	old, _ := s.ReadOne(&t.Sku.Kennung)

	if old == nil || old.Less(*t) {
		s.delegate.addOne(t)
	}

	if t.IsNew() {
		s.LogWriter.New(t)
	} else {
		s.LogWriter.Updated(t)
	}

	return
}

func (s *commonStoreBase[O, OPtr, K, KPtr]) GetInheritor(
	orf schnittstellen.ObjekteReaderFactory,
	arf schnittstellen.AkteReaderFactory,
	pmf objekte_format.Format,
) objekte_store.TransactedInheritor {
	p := collections.MakePool[
		objekte.Transacted[O, OPtr, K, KPtr],
		*objekte.Transacted[O, OPtr, K, KPtr],
	]()

	inflator := objekte_store.MakeTransactedInflator[
		O,
		OPtr,
		K,
		KPtr,
	](
		s.StoreUtil.GetStoreVersion(),
		schnittstellen.MakeBespokeObjekteReadWriterFactory(orf, s),
		schnittstellen.MakeBespokeAkteReadWriterFactory(arf, s),
		pmf,
		objekte_store.MakeAkteFormat[O, OPtr](
			objekte.MakeNopAkteParseSaver[O, OPtr](s),
			nil,
			s,
		),
		p,
	)

	return objekte_store.MakeTransactedInheritor[
		objekte.Transacted[O, OPtr, K, KPtr],
		*objekte.Transacted[O, OPtr, K, KPtr],
	](
		inflator,
		s,
		p,
	)
}
