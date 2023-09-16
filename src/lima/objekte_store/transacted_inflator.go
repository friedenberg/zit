package objekte_store

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/juliett/objekte"
)

type TransactedDataIdentityInflator[A any] interface {
	InflateFromSku(sku.SkuLike) (A, error)
}

type ObjekteStorer[A any] interface {
	StoreObjekte(A) error
}

type AkteStorer[A any] interface {
	StoreAkte(A) error
}

// TODO-P1 split into ObjekteInflator
type TransactedInflator[
	A objekte.Akte[A],
	APtr objekte.AktePtr[A],
	K kennung.KennungLike[K],
	KPtr kennung.KennungLikePtr[K],
] interface {
	InflateFromSkuLike(
		sku.SkuLike,
	) (*sku.Transacted[K, KPtr], error)
	InflatorStorer[*sku.Transacted[K, KPtr]]
	InflateFromSkuAndStore(sku.SkuLike) error
}

type transactedInflator[
	A objekte.Akte[A],
	APtr objekte.AktePtr[A],
	K kennung.KennungLike[K],
	KPtr kennung.KennungLikePtr[K],
] struct {
	storeVersion              schnittstellen.StoreVersion
	of                        schnittstellen.ObjekteIOFactory
	af                        schnittstellen.AkteIOFactory
	persistentMetadateiFormat objekte_format.Format
	options                   objekte_format.Options
	akteFormat                objekte.AkteFormat[A, APtr]
	pool                      schnittstellen.Pool[
		sku.Transacted[K, KPtr],
		*sku.Transacted[K, KPtr],
	]
}

func MakeTransactedInflator[
	A objekte.Akte[A],
	APtr objekte.AktePtr[A],
	K kennung.KennungLike[K],
	KPtr kennung.KennungLikePtr[K],
](
	sv schnittstellen.StoreVersion,
	of schnittstellen.ObjekteIOFactory,
	af schnittstellen.AkteIOFactory,
	persistentMetadateiFormat objekte_format.Format,
	op objekte_format.Options,
	akteFormat objekte.AkteFormat[A, APtr],
	pool schnittstellen.Pool[
		sku.Transacted[K, KPtr],
		*sku.Transacted[K, KPtr],
	],
) *transactedInflator[A, APtr, K, KPtr] {
	return &transactedInflator[A, APtr, K, KPtr]{
		storeVersion:              sv,
		of:                        of,
		af:                        af,
		persistentMetadateiFormat: persistentMetadateiFormat,
		options:                   op,
		akteFormat:                akteFormat,
		pool:                      pool,
	}
}

func (h *transactedInflator[A, APtr, K, KPtr]) InflateFromSkuLike(
	o sku.SkuLike,
) (t *sku.Transacted[K, KPtr], err error) {
	if h.pool == nil {
		t = new(sku.Transacted[K, KPtr])
	} else {
		t = h.pool.Get()
	}

	// TODO-P2 make generic
	if err = t.SetFromSkuLike(o); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.SetTai(o.GetTai())

	// TODO-P2 make generic
	if t.GetGattung() != o.GetGattung() {
		err = errors.Errorf(
			"expected gattung %s but got %s",
			t.GetGattung(),
			o.GetGattung(),
		)
		return
	}

	if err = KPtr(&t.Kennung).Set(o.GetKennungLike().String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.ObjekteSha = sha.Make(o.GetObjekteSha())

	if h.storeVersion.GetInt() < 3 {
		if err = h.readObjekte(o, t); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	// TODO-P1 switch to pool
	var a1 A
	a := APtr(&a1)

	if err = h.readAkte(t, a); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (h *transactedInflator[A, APtr, K, KPtr]) InflateFromSku(
	o sku.SkuLike,
) (t *sku.Transacted[K, KPtr], err error) {
	if h.pool == nil {
		t = new(sku.Transacted[K, KPtr])
	} else {
		t = h.pool.Get()
	}

	if err = t.SetFromSkuLike(o); err != nil {
		err = errors.Wrapf(err, "Sku: %s", o)
		return
	}

	t.GetTai()

	if h.storeVersion.GetInt() < 3 {
		if err = h.readObjekte(o, t); err != nil {
			err = errors.Wrapf(err, "Sku: %s", o)
			return
		}
	}

	// TODO-P1 switch to pool
	var a1 A
	a := APtr(&a1)

	if err = h.readAkte(t, a); err != nil {
		err = errors.Wrapf(err, "Sku: %s", o)
		return
	}

	return
}

func (h *transactedInflator[A, APtr, K, KPtr]) StoreAkte(
	t *sku.Transacted[K, KPtr],
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

func (h *transactedInflator[A, APtr, K, KPtr]) StoreObjekte(
	t *sku.Transacted[K, KPtr],
) (err error) {
	if h.storeVersion.GetInt() >= 3 {
		return
	}

	var ow sha.WriteCloser

	if ow, err = h.of.ObjekteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, ow)

	if _, err = h.persistentMetadateiFormat.FormatPersistentMetadatei(
		ow,
		t,
		h.options,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.ObjekteSha = sha.Make(ow.GetShaLike())
	t.SetAkteSha(t.GetAkteSha())

	return
}

func (h *transactedInflator[A, APtr, K, KPtr]) InflateFromSkuAndStore(
	o sku.SkuLike,
) (err error) {
	var t *sku.Transacted[K, KPtr]

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

func (h *transactedInflator[A, APtr, K, KPtr]) readObjekte(
	sk sku.SkuLike,
	t *sku.Transacted[K, KPtr],
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
		h.options,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.ObjekteSha = sha.Make(r.GetShaLike())

	if !t.ObjekteSha.EqualsSha(sk.GetObjekteSha()) {
		errors.Todo(
			"objekte sha mismatch for %s! expected %s but got %s.",
			sk.GetGattung(),
			sk.GetObjekteSha(),
			t.ObjekteSha,
		)
	}

	errors.Log().Printf("parsed %d objekte bytes", n)

	return
}

func (h *transactedInflator[A, APtr, K, KPtr]) readAkte(
	t *sku.Transacted[K, KPtr],
	a APtr,
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

	var n int64

	sw := sha.MakeWriter(io.Discard)

	if n, err = h.akteFormat.ParseAkte(io.TeeReader(r, sw), a); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh := sw.GetShaLike()

	if !t.GetAkteSha().EqualsSha(sh) {
		errors.TodoRecoverable(
			"objekte had akte sha %s, but akte reader had sha %s +%d",
			t.GetAkteSha(),
			sh,
			n,
		)
	}

	t.SetAkteSha(sh)

	errors.Log().Printf("parsed %d akte bytes: %s", n, sh)

	return
}
