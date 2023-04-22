package objekte_store

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/golf/persisted_metadatei_format"
)

type TransactedDataIdentityInflator[T any] interface {
	InflateFromDataIdentity(sku.DataIdentity) (T, error)
}

type ObjekteStorer[T any] interface {
	StoreObjekte(T) error
}

type AkteStorer[T any] interface {
	StoreAkte(T) error
}

// TODO-P1 split into ObjekteInflator
type TransactedInflator[
	T objekte.Objekte[T],
	T1 objekte.ObjektePtr[T],
	T2 schnittstellen.Id[T2],
	T3 schnittstellen.IdPtr[T2],
	T4 any,
	T5 schnittstellen.VerzeichnissePtr[T4, T],
] interface {
	InflateFromSku(sku.Sku) (*objekte.Transacted[T, T1, T2, T3, T4, T5], error)
	InflateFromSku2(sku.Sku2) (*objekte.Transacted[T, T1, T2, T3, T4, T5], error)
	InflatorStorer[*objekte.Transacted[T, T1, T2, T3, T4, T5]]
	InflateFromDataIdentityAndStore(sku.DataIdentity) error
}

type transactedInflator[
	T objekte.Objekte[T],
	T1 objekte.ObjektePtr[T],
	T2 schnittstellen.Id[T2],
	T3 schnittstellen.IdPtr[T2],
	T4 any,
	T5 schnittstellen.VerzeichnissePtr[T4, T],
] struct {
	of                        schnittstellen.ObjekteIOFactory
	af                        schnittstellen.AkteIOFactory
	persistentMetadateiFormat persisted_metadatei_format.V0
	akteFormat                objekte.AkteFormat[T, T1]
	pool                      schnittstellen.Pool[
		objekte.Transacted[T, T1, T2, T3, T4, T5],
		*objekte.Transacted[T, T1, T2, T3, T4, T5],
	]
}

func MakeTransactedInflator[
	T objekte.Objekte[T],
	T1 objekte.ObjektePtr[T],
	T2 schnittstellen.Id[T2],
	T3 schnittstellen.IdPtr[T2],
	T4 any,
	T5 schnittstellen.VerzeichnissePtr[T4, T],
](
	of schnittstellen.ObjekteIOFactory,
	af schnittstellen.AkteIOFactory,
	persistentMetadateiFormat persisted_metadatei_format.V0,
	akteFormat objekte.AkteFormat[T, T1],
	pool schnittstellen.Pool[
		objekte.Transacted[T, T1, T2, T3, T4, T5],
		*objekte.Transacted[T, T1, T2, T3, T4, T5],
	],
) *transactedInflator[T, T1, T2, T3, T4, T5] {
	return &transactedInflator[T, T1, T2, T3, T4, T5]{
		of:                        of,
		af:                        af,
		persistentMetadateiFormat: persistentMetadateiFormat,
		akteFormat:                akteFormat,
		pool:                      pool,
	}
}

func (h *transactedInflator[T, T1, T2, T3, T4, T5]) InflateFromSku2(
	o sku.Sku2,
) (t *objekte.Transacted[T, T1, T2, T3, T4, T5], err error) {
	if h.pool == nil {
		t = new(objekte.Transacted[T, T1, T2, T3, T4, T5])
	} else {
		t = h.pool.Get()
	}

	// TODO make generic
	if err = t.Sku.SetFromSku2(o); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO make generic
	if t.Sku.Kennung.GetGattung() != o.Gattung {
		err = errors.Errorf(
			"expected gattung %s but got %s",
			t.Sku.Kennung.GetGattung(),
			o.Gattung,
		)
		return
	}

	if err = T3(&t.Sku.Kennung).Set(o.Kennung.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.Sku.ObjekteSha = o.ObjekteSha

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

func (h *transactedInflator[T, T1, T2, T3, T4, T5]) InflateFromSku(
	o sku.Sku,
) (t *objekte.Transacted[T, T1, T2, T3, T4, T5], err error) {
	if h.pool == nil {
		t = new(objekte.Transacted[T, T1, T2, T3, T4, T5])
	} else {
		t = h.pool.Get()
	}

	// TODO make generic
	if err = t.Sku.SetFromSku(o); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO make generic
	if t.Sku.Kennung.GetGattung() != o.Gattung {
		err = errors.Errorf(
			"expected gattung %s but got %s",
			t.Sku.Kennung.GetGattung(),
			o.Gattung,
		)
		return
	}

	if err = T3(&t.Sku.Kennung).Set(o.Kennung.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.Sku.ObjekteSha = o.ObjekteSha

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

func (h *transactedInflator[T, T1, T2, T3, T4, T5]) InflateFromDataIdentity(
	o sku.DataIdentity,
) (t *objekte.Transacted[T, T1, T2, T3, T4, T5], err error) {
	if h.pool == nil {
		t = new(objekte.Transacted[T, T1, T2, T3, T4, T5])
	} else {
		t = h.pool.Get()
	}

	if err = t.SetDataIdentity(o); err != nil {
		err = errors.Wrapf(err, "DataIdentity: %s", o)
		return
	}

	if err = h.readObjekte(o, t); err != nil {
		err = errors.Wrapf(err, "DataIdentity: %s", o)
		return
	}

	if err = h.readAkte(t); err != nil {
		err = errors.Wrapf(err, "DataIdentity: %s", o)
		return
	}

	return
}

func (h *transactedInflator[T, T1, T2, T3, T4, T5]) StoreAkte(
	t *objekte.Transacted[T, T1, T2, T3, T4, T5],
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

	t.SetAkteSha(aw.Sha())
	t.Sku.AkteSha = sha.Make(aw.Sha())

	return
}

func (h *transactedInflator[T, T1, T2, T3, T4, T5]) StoreObjekte(
	t *objekte.Transacted[T, T1, T2, T3, T4, T5],
) (err error) {
	var ow sha.WriteCloser

	if ow, err = h.of.ObjekteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, ow)

	if _, err = h.persistentMetadateiFormat.Format(
		ow,
		t.Objekte,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.Sku.ObjekteSha = sha.Make(ow.Sha())
	t.Sku.AkteSha = sha.Make(t.GetAkteSha())

	return
}

func (h *transactedInflator[T, T1, T2, T3, T4, T5]) InflateFromDataIdentityAndStore(
	o sku.DataIdentity,
) (err error) {
	var t *objekte.Transacted[T, T1, T2, T3, T4, T5]

	if t, err = h.InflateFromDataIdentity(o); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = h.StoreObjekte(t); err != nil {
		err = errors.Wrapf(err, "DataIdentity: %s", o)
		return
	}

	return
}

func (h *transactedInflator[T, T1, T2, T3, T4, T5]) readObjekte(
	sk sku.DataIdentity,
	t *objekte.Transacted[T, T1, T2, T3, T4, T5],
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

	if n, err = h.persistentMetadateiFormat.Parse(r, T1(&t.Objekte)); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.Sku.ObjekteSha = sha.Make(r.Sha())

	if !t.Sku.ObjekteSha.EqualsSha(sk.GetObjekteSha()) {
		errors.Todo(
			"objekte sha mismatch for %s! expected %s but got %s.\nObjekte: %v",
			sk.GetGattung(),
			sk.GetObjekteSha(),
			t.Sku.ObjekteSha,
			t.Objekte,
		)
	}

	T5(&t.Verzeichnisse).ResetWithObjekte(t.Objekte)

	errors.Log().Printf("parsed %d objekte bytes", n)

	return
}

func (h *transactedInflator[T, T1, T2, T3, T4, T5]) readAkte(
	t *objekte.Transacted[T, T1, T2, T3, T4, T5],
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
		sh schnittstellen.Sha
	)

	if sh, n, err = h.akteFormat.ParseSaveAkte(r, &t.Objekte); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.SetAkteSha(sh)

	errors.Log().Printf("parsed %d akte bytes: %s", n, sh)

	return
}
