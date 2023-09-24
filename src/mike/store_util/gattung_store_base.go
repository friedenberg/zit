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
] struct {
	schnittstellen.GattungGetter

	schnittstellen.ObjekteIOFactory

	delegate CommonStoreDelegate

	StoreUtil

	Pool schnittstellen.Pool[
		sku.Transacted,
		*sku.Transacted,
	]

	objekte_store.TransactedInflator

	objekte_store.AkteTextSaver[O, OPtr]

	objekte_store.StoredParseSaver[O, OPtr]

	objekte_store.TransactedReader

	objekte_store.LogWriter[*sku.Transacted]

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
	delegate CommonStoreDelegate,
	sa StoreUtil,
	tr objekte_store.TransactedReader,
	pmf objekte_format.Format,
	akteFormat objekte.AkteFormat[O, OPtr],
) (s *CommonStoreBase[O, OPtr], err error) {
	// type T objekte.Transacted[O, OPtr, K, KPtr, ]
	// type TPtr *objekte.Transacted[O, OPtr, K, KPtr, ]

	if delegate == nil {
		panic("delegate was nil")
	}

	pool := pool.MakePoolWithReset[
		sku.Transacted,
		*sku.Transacted,
	]()

	of := sa.ObjekteReaderWriterFactory(gg)

	s = &CommonStoreBase[O, OPtr]{
		GattungGetter:    gg,
		ObjekteIOFactory: of,
		delegate:         delegate,
		StoreUtil:        sa,
		Pool:             pool,
		akteFormat:       akteFormat,
		TransactedInflator: objekte_store.MakeTransactedInflator[
			O,
			OPtr,
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

func (s *CommonStoreBase[O, OPtr]) SetLogWriter(
	lw objekte_store.LogWriter[*sku.Transacted],
) {
	s.LogWriter = lw
}

func (s *CommonStoreBase[O, OPtr]) Query(
	m matcher.MatcherSigil,
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	return objekte_store.QueryMethodForMatcher(s, m, f)
}

func (s *CommonStoreBase[O, OPtr]) ReindexOne(
	sk *sku.Transacted,
) (o matcher.Matchable, err error) {
	var t *sku.Transacted

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

func (s *CommonStoreBase[O, OPtr]) Inherit(
	t *sku.Transacted,
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

func (s *CommonStoreBase[O, OPtr]) GetInheritor(
	orf schnittstellen.ObjekteReaderFactory,
	arf schnittstellen.AkteReaderFactory,
	pmf objekte_format.Format,
) objekte_store.TransactedInheritor {
	p := pool.MakePoolWithReset[
		sku.Transacted,
		*sku.Transacted,
	]()

	inflator := objekte_store.MakeTransactedInflator[
		O,
		OPtr,
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
		sku.Transacted,
		*sku.Transacted,
	](
		inflator,
		s,
		p,
	)
}

func (s *CommonStoreBase[O, OPtr]) GetAkte(
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

func (s *CommonStoreBase[O, OPtr]) PutAkte(a OPtr) {
	// TODO-P2 implement pool
}
