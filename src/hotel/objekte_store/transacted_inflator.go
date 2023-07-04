package objekte_store

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/golf/objekte_format"
)

type TransactedDataIdentityInflator[T any] interface {
	InflateFromSku(sku.SkuLike) (T, error)
}

type ObjekteStorer[T any] interface {
	StoreObjekte(T) error
}

type AkteStorer[T any] interface {
	StoreAkte(T) error
}

// TODO-P1 split into ObjekteInflator
type TransactedInflator[
	T objekte.Akte[T],
	T1 objekte.AktePtr[T],
	T2 kennung.KennungLike[T2],
	T3 kennung.KennungLikePtr[T2],
] interface {
	InflateFromSkuLike(sku.SkuLike) (*objekte.Transacted[T, T1, T2, T3], error)
	InflatorStorer[*objekte.Transacted[T, T1, T2, T3]]
	InflateFromSkuAndStore(sku.SkuLike) error
}

type transactedInflator[
	T objekte.Akte[T],
	T1 objekte.AktePtr[T],
	T2 kennung.KennungLike[T2],
	T3 kennung.KennungLikePtr[T2],
] struct {
	of                        schnittstellen.ObjekteIOFactory
	af                        schnittstellen.AkteIOFactory
	persistentMetadateiFormat objekte_format.Format
	akteFormat                objekte.AkteFormat[T, T1]
	pool                      schnittstellen.Pool[
		objekte.Transacted[T, T1, T2, T3],
		*objekte.Transacted[T, T1, T2, T3],
	]
}

func MakeTransactedInflator[
	T objekte.Akte[T],
	T1 objekte.AktePtr[T],
	T2 kennung.KennungLike[T2],
	T3 kennung.KennungLikePtr[T2],
](
	of schnittstellen.ObjekteIOFactory,
	af schnittstellen.AkteIOFactory,
	persistentMetadateiFormat objekte_format.Format,
	akteFormat objekte.AkteFormat[T, T1],
	pool schnittstellen.Pool[
		objekte.Transacted[T, T1, T2, T3],
		*objekte.Transacted[T, T1, T2, T3],
	],
) *transactedInflator[T, T1, T2, T3] {
	return &transactedInflator[T, T1, T2, T3]{
		of:                        of,
		af:                        af,
		persistentMetadateiFormat: persistentMetadateiFormat,
		akteFormat:                akteFormat,
		pool:                      pool,
	}
}

func (h *transactedInflator[T, T1, T2, T3]) InflateFromSkuLike(
	o sku.SkuLike,
) (t *objekte.Transacted[T, T1, T2, T3], err error) {
	if h.pool == nil {
		t = new(objekte.Transacted[T, T1, T2, T3])
	} else {
		t = h.pool.Get()
	}

	// TODO-P2 make generic
	if err = t.Sku.SetFromSkuLike(o); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.SetTai(o.GetTai())

	// TODO-P2 make generic
	if t.Sku.GetGattung() != o.GetGattung() {
		err = errors.Errorf(
			"expected gattung %s but got %s",
			t.Sku.GetGattung(),
			o.GetGattung(),
		)
		return
	}

	if err = T3(&t.Sku.Kennung).Set(o.GetKennungLike().String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.Sku.ObjekteSha = sha.Make(o.GetObjekteSha())

	if err = h.readObjekte(o, t); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = h.readAkte(t); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (h *transactedInflator[T, T1, T2, T3]) InflateFromSku(
	o sku.SkuLike,
) (t *objekte.Transacted[T, T1, T2, T3], err error) {
	if h.pool == nil {
		t = new(objekte.Transacted[T, T1, T2, T3])
	} else {
		t = h.pool.Get()
	}

	if err = t.Sku.SetFromSkuLike(o); err != nil {
		err = errors.Wrapf(err, "Sku: %s", o)
		return
	}

	t.GetTai()

	if err = h.readObjekte(o, t); err != nil {
		err = errors.Wrapf(err, "Sku: %s", o)
		return
	}

	if err = h.readAkte(t); err != nil {
		err = errors.Wrapf(err, "Sku: %s", o)
		return
	}

	return
}

func (h *transactedInflator[T, T1, T2, T3]) StoreAkte(
	t *objekte.Transacted[T, T1, T2, T3],
) (err error) {
	var aw sha.WriteCloser

	if aw, err = h.af.AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, aw)

	if _, err = h.akteFormat.FormatSavedAkte(aw, t.GetAkteSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.SetAkteSha(aw.GetShaLike())

	return
}

func (h *transactedInflator[T, T1, T2, T3]) StoreObjekte(
	t *objekte.Transacted[T, T1, T2, T3],
) (err error) {
	var ow sha.WriteCloser

	if ow, err = h.of.ObjekteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, ow)

	if _, err = h.persistentMetadateiFormat.FormatPersistentMetadatei(
		ow,
		t,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.Sku.ObjekteSha = sha.Make(ow.GetShaLike())
	t.SetAkteSha(t.GetAkteSha())

	return
}

func (h *transactedInflator[T, T1, T2, T3]) InflateFromSkuAndStore(
	o sku.SkuLike,
) (err error) {
	var t *objekte.Transacted[T, T1, T2, T3]

	if t, err = h.InflateFromSku(o); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = h.StoreObjekte(t); err != nil {
		err = errors.Wrapf(err, "Sku: %s", o)
		return
	}

	return
}

func (h *transactedInflator[T, T1, T2, T3]) readObjekte(
	sk sku.SkuLike,
	t *objekte.Transacted[T, T1, T2, T3],
) (err error) {
	if sk.GetObjekteSha().IsNull() {
		return
	}

	var r sha.ReadCloser

	if r, err = h.of.ObjekteReader(sk.GetObjekteSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, r)

	var n int64

	if n, err = h.persistentMetadateiFormat.ParsePersistentMetadatei(
		r,
		t,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.Sku.ObjekteSha = sha.Make(r.GetShaLike())

	if !t.Sku.ObjekteSha.EqualsSha(sk.GetObjekteSha()) {
		errors.Todo(
			"objekte sha mismatch for %s! expected %s but got %s.\nObjekte: %v",
			sk.GetGattung(),
			sk.GetObjekteSha(),
			t.Sku.ObjekteSha,
			t.Akte,
		)
	}

	errors.Log().Printf("parsed %d objekte bytes", n)

	return
}

func (h *transactedInflator[T, T1, T2, T3]) readAkte(
	t *objekte.Transacted[T, T1, T2, T3],
) (err error) {
	if h.akteFormat == nil {
		return
	}

	if t.GetAkteSha().IsNull() {
		return
	}

	var r sha.ReadCloser

	if r, err = h.af.AkteReader(t.GetAkteSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, r)

	var (
		n  int64
		sh schnittstellen.ShaLike
	)

	if sh, n, err = h.akteFormat.ParseSaveAkte(r, &t.Akte); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.SetAkteSha(sh)

	errors.Log().Printf("parsed %d akte bytes: %s", n, sh)

	return
}
