package store_util

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/alfa/toml"
	"github.com/friedenberg/zit/src/charlie/pool"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/juliett/objekte"
	"github.com/friedenberg/zit/src/lima/objekte_store"
)

type GattungStoreLike interface {
	Reindexer
	schnittstellen.ObjekteIOFactory
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
] struct {
	schnittstellen.GattungGetter

	schnittstellen.ObjekteIOFactory

	delegate CommonStoreDelegate[O, OPtr, K, KPtr]

	StoreUtil

	Pool schnittstellen.Pool[
		sku.Transacted[K, KPtr],
		*sku.Transacted[K, KPtr],
	]

	objekte_store.TransactedInflator[O, OPtr, K, KPtr]

	objekte_store.AkteTextSaver[O, OPtr]

	objekte_store.StoredParseSaver[O, OPtr, K, KPtr]

	objekte_store.TransactedReader[
		KPtr,
		sku.SkuLikePtr,
	]

	objekte_store.LogWriter[sku.SkuLikePtr]

	persistentMetadateiFormat objekte_format.Format

	akteFormat objekte.AkteFormat[O, OPtr]
}

func MakeCommonStoreBase[
	O objekte.Akte[O],
	OPtr objekte.AktePtr[O],
	K kennung.KennungLike[K],
	KPtr kennung.KennungLikePtr[K],
](
	gg schnittstellen.GattungGetter,
	delegate CommonStoreDelegate[O, OPtr, K, KPtr],
	sa StoreUtil,
	tr objekte_store.TransactedReader[
		KPtr,
		sku.SkuLikePtr,
	],
	pmf objekte_format.Format,
	akteFormat objekte.AkteFormat[O, OPtr],
) (s *CommonStoreBase[O, OPtr, K, KPtr], err error) {
	// type T objekte.Transacted[O, OPtr, K, KPtr, ]
	// type TPtr *objekte.Transacted[O, OPtr, K, KPtr, ]

	if delegate == nil {
		panic("delegate was nil")
	}

	pool := pool.MakePoolWithReset[
		sku.Transacted[K, KPtr],
		*sku.Transacted[K, KPtr],
	]()

	of := sa.ObjekteReaderWriterFactory(gg)

	s = &CommonStoreBase[O, OPtr, K, KPtr]{
		GattungGetter:    gg,
		ObjekteIOFactory: of,
		delegate:         delegate,
		StoreUtil:        sa,
		Pool:             pool,
		akteFormat:       akteFormat,
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
			objekte_format.Options{IncludeTai: true},
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

func (s *CommonStoreBase[O, OPtr, K, KPtr]) SetLogWriter(
	lw objekte_store.LogWriter[sku.SkuLikePtr],
) {
	s.LogWriter = lw
}

func (s *CommonStoreBase[O, OPtr, K, KPtr]) Query(
	m matcher.MatcherSigil,
	f schnittstellen.FuncIter[sku.SkuLikePtr],
) (err error) {
	return objekte_store.QueryMethodForMatcher[
		KPtr,
		sku.SkuLikePtr,
	](s, m, f)
}

func (s *CommonStoreBase[O, OPtr, K, KPtr]) ReindexOne(
	sk sku.SkuLike,
) (o matcher.Matchable, err error) {
	var t *sku.Transacted[K, KPtr]

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
		if err = s.delegate.AddOne(t); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		s.LogWriter.Updated(t)
		if err = s.delegate.UpdateOne(t); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s *CommonStoreBase[O, OPtr, K, KPtr]) Inherit(
	t *sku.Transacted[K, KPtr],
) (err error) {
	if t == nil {
		err = errors.Errorf("trying to inherit nil %T", t)
		return
	}

	errors.Log().Printf("inheriting %s", t.ObjekteSha)

	s.StoreUtil.AddMatchable(t)
	s.StoreUtil.CommitTransacted(t)

	old, _ := s.ReadOne(t.GetKennungLike())

	if old == nil || old.GetTai().Less(t.GetTai()) {
		s.delegate.AddOne(t)
	}

	if t.IsNew() {
		s.LogWriter.New(t)
	} else {
		s.LogWriter.Updated(t)
	}

	return
}

func (s *CommonStoreBase[O, OPtr, K, KPtr]) GetInheritor(
	orf schnittstellen.ObjekteReaderFactory,
	arf schnittstellen.AkteReaderFactory,
	pmf objekte_format.Format,
) objekte_store.TransactedInheritor {
	p := pool.MakePoolWithReset[
		sku.Transacted[K, KPtr],
		*sku.Transacted[K, KPtr],
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
		objekte_format.Options{IncludeTai: true},
		objekte_store.MakeAkteFormat[O, OPtr](
			objekte.MakeNopAkteParseSaver[O, OPtr](s),
			nil,
			s,
		),
		p,
	)

	return objekte_store.MakeTransactedInheritor[
		sku.Transacted[K, KPtr],
		*sku.Transacted[K, KPtr],
	](
		inflator,
		s,
		p,
	)
}

func (s *CommonStoreBase[O, OPtr, K, KPtr]) GetAkte(
	sh schnittstellen.ShaLike,
) (a OPtr, err error) {
	var ar schnittstellen.ShaReadCloser

	if ar, err = s.StoreUtil.AkteReader(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, ar)

	var a1 O
	a = OPtr(&a1)
	a.Reset()

	if _, err = s.akteFormat.ParseAkte(ar, a); err != nil {
		err = errors.Wrap(err)
		return
	}

	actual := ar.GetShaLike()

	if !actual.EqualsSha(sh) {
		err = errors.Errorf("expected sha %s but got %s", sh, actual)
		return
	}

	return
}

func (s *CommonStoreBase[O, OPtr, K, KPtr]) PutAkte(a OPtr) {
	// TODO-P2 implement pool
}