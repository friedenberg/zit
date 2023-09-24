package objekte_store

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/juliett/objekte"
)

type TransactedDataIdentityInflator interface {
	InflateFromSku(*sku.Transacted) (*sku.Transacted, error)
}

type AkteStorer[A any] interface {
	StoreAkte(A) error
}

// TODO-P1 split into ObjekteInflator
type TransactedInflator interface {
	InflateFromSkuLike(
		*sku.Transacted,
	) (*sku.Transacted, error)
	InflatorStorer
	InflateFromSkuAndStore(*sku.Transacted) error
}

type transactedInflator[
	A objekte.Akte[A],
	APtr objekte.AktePtr[A],
] struct {
	storeVersion              schnittstellen.StoreVersion
	af                        schnittstellen.AkteIOFactory
	persistentMetadateiFormat objekte_format.Format
	options                   objekte_format.Options
	akteFormat                objekte.AkteFormat[A, APtr]
}

func MakeTransactedInflator[
	A objekte.Akte[A],
	APtr objekte.AktePtr[A],
](
	sv schnittstellen.StoreVersion,
	af schnittstellen.AkteIOFactory,
	persistentMetadateiFormat objekte_format.Format,
	op objekte_format.Options,
	akteFormat objekte.AkteFormat[A, APtr],
) *transactedInflator[A, APtr] {
	return &transactedInflator[A, APtr]{
		storeVersion:              sv,
		af:                        af,
		persistentMetadateiFormat: persistentMetadateiFormat,
		options:                   op,
		akteFormat:                akteFormat,
	}
}

func (h *transactedInflator[A, APtr]) InflateFromSkuLike(
	o *sku.Transacted,
) (t *sku.Transacted, err error) {
	t = sku.GetTransactedPool().Get()

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

	if err = t.Kennung.SetWithKennung(o.GetKennungLike()); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.ObjekteSha = sha.Make(o.GetObjekteSha())

	// TODO-P1 switch to pool
	var a1 A
	a := APtr(&a1)

	if err = h.readAkte(t, a); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (h *transactedInflator[A, APtr]) InflateFromSku(
	o *sku.Transacted,
) (t *sku.Transacted, err error) {
	t = sku.GetTransactedPool().Get()

	if err = t.SetFromSkuLike(o); err != nil {
		err = errors.Wrapf(err, "Sku: %s", o)
		return
	}

	t.GetTai()

	// TODO-P1 switch to pool
	var a1 A
	a := APtr(&a1)

	if err = h.readAkte(t, a); err != nil {
		err = errors.Wrapf(err, "Sku: %s", o)
		return
	}

	return
}

func (h *transactedInflator[A, APtr]) StoreAkte(
	t *sku.Transacted,
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

func (h *transactedInflator[A, APtr]) InflateFromSkuAndStore(
	o *sku.Transacted,
) (err error) {
	if _, err = h.InflateFromSku(o); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (h *transactedInflator[A, APtr]) readAkte(
	t *sku.Transacted,
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
