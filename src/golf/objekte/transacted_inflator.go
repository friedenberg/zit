package objekte

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/foxtrot/sku"
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
	T schnittstellen.Objekte[T],
	T1 schnittstellen.ObjektePtr[T],
	T2 schnittstellen.Id[T2],
	T3 schnittstellen.IdPtr[T2],
	T4 any,
	T5 schnittstellen.VerzeichnissePtr[T4, T],
] interface {
	InflateFromSku(sku.Sku) (*Transacted[T, T1, T2, T3, T4, T5], error)
	InflateFromSku2(sku.Sku2) (*Transacted[T, T1, T2, T3, T4, T5], error)
	InflatorStorer[*Transacted[T, T1, T2, T3, T4, T5]]
	InflateFromDataIdentityAndStore(sku.DataIdentity) error
}

type transactedInflator[
	T schnittstellen.Objekte[T],
	T1 schnittstellen.ObjektePtr[T],
	T2 schnittstellen.Id[T2],
	T3 schnittstellen.IdPtr[T2],
	T4 any,
	T5 schnittstellen.VerzeichnissePtr[T4, T],
] struct {
	of            schnittstellen.ObjekteIOFactory
	af            schnittstellen.AkteIOFactory
	objekteFormat schnittstellen.Format[T, T1]
	akteFormat    schnittstellen.Format[T, T1]
	pool          collections.PoolLike[Transacted[T, T1, T2, T3, T4, T5]]
}

func MakeTransactedInflator[
	T schnittstellen.Objekte[T],
	T1 schnittstellen.ObjektePtr[T],
	T2 schnittstellen.Id[T2],
	T3 schnittstellen.IdPtr[T2],
	T4 any,
	T5 schnittstellen.VerzeichnissePtr[T4, T],
](
	of schnittstellen.ObjekteIOFactory,
	af schnittstellen.AkteIOFactory,
	objekteFormat schnittstellen.Format[T, T1],
	akteFormat schnittstellen.Format[T, T1],
	pool collections.PoolLike[Transacted[T, T1, T2, T3, T4, T5]],
) *transactedInflator[T, T1, T2, T3, T4, T5] {
	if objekteFormat == nil {
		objekteFormat = MakeFormat[T, T1]()
	}

	return &transactedInflator[T, T1, T2, T3, T4, T5]{
		of:            of,
		af:            af,
		objekteFormat: objekteFormat,
		akteFormat:    akteFormat,
		pool:          pool,
	}
}

func (h *transactedInflator[T, T1, T2, T3, T4, T5]) InflateFromSku2(
	o sku.Sku2,
) (t *Transacted[T, T1, T2, T3, T4, T5], err error) {
	if h.pool == nil {
		t = new(Transacted[T, T1, T2, T3, T4, T5])
	} else {
		t = h.pool.Get()
	}

	//TODO make generic
	if err = t.Sku.SetFromSku2(o); err != nil {
		err = errors.Wrap(err)
		return
	}

	//TODO make generic
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
) (t *Transacted[T, T1, T2, T3, T4, T5], err error) {
	if h.pool == nil {
		t = new(Transacted[T, T1, T2, T3, T4, T5])
	} else {
		t = h.pool.Get()
	}

	//TODO make generic
	if err = t.Sku.SetFromSku(o); err != nil {
		err = errors.Wrap(err)
		return
	}

	//TODO make generic
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
) (t *Transacted[T, T1, T2, T3, T4, T5], err error) {
	if h.pool == nil {
		t = new(Transacted[T, T1, T2, T3, T4, T5])
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
	t *Transacted[T, T1, T2, T3, T4, T5],
) (err error) {
	var aw sha.WriteCloser

	if aw, err = h.af.AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, aw)

	if _, err = h.akteFormat.Format(aw, &t.Objekte); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.SetAkteSha(aw.Sha())
	t.Sku.AkteSha = sha.Make(aw.Sha())

	return
}

func (h *transactedInflator[T, T1, T2, T3, T4, T5]) StoreObjekte(
	t *Transacted[T, T1, T2, T3, T4, T5],
) (err error) {
	var ow sha.WriteCloser

	if ow, err = h.of.ObjekteWriter(t.GetGattung()); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, ow)

	if _, err = h.objekteFormat.Format(ow, &t.Objekte); err != nil {
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
	var t *Transacted[T, T1, T2, T3, T4, T5]

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
	t *Transacted[T, T1, T2, T3, T4, T5],
) (err error) {
	if sk.GetObjekteSha().IsNull() {
		return
	}

	var r sha.ReadCloser

	if r, err = h.of.ObjekteReader(sk, sk.GetObjekteSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, r)

	var n int64

	if n, err = h.objekteFormat.Parse(r, &t.Objekte); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.Sku.ObjekteSha = sha.Make(r.Sha())

	if !t.Sku.ObjekteSha.Equals(sk.GetObjekteSha()) {
		errors.Todo(
			"objekte sha mismatch for %s! expected %s but got %s.\nObjekte: %v",
			sk.GetGattung(),
			sk.GetObjekteSha(),
			t.Sku.ObjekteSha,
			t.Objekte,
		)
	}

	errors.Log().Printf("parsed %d objekte bytes", n)

	return
}

func (h *transactedInflator[T, T1, T2, T3, T4, T5]) readAkte(
	t *Transacted[T, T1, T2, T3, T4, T5],
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

	if n, err = h.akteFormat.Parse(r, &t.Objekte); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.Log().Printf("parsed %d akte bytes", n)

	return
}