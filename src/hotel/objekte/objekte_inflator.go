package objekte

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/schnittstellen"
)

type ObjekteInflator[
	T schnittstellen.Objekte[T],
	T1 schnittstellen.ObjektePtr[T],
	T4 any,
	T5 schnittstellen.VerzeichnissePtr[T4, T],
] interface {
	InflateObjekteFromSku(sku.Sku) (T1, error)
}

type objekteInflator[
	T schnittstellen.Objekte[T],
	T1 schnittstellen.ObjektePtr[T],
	T2 any,
	T3 schnittstellen.VerzeichnissePtr[T2, T],
] struct {
	or            schnittstellen.ObjekteReaderFactory
	ar            schnittstellen.AkteReaderFactory
	objekteParser schnittstellen.Parser[T, T1]
	akteParser    schnittstellen.Parser[T, T1]
	pool          collections.PoolLike[T]
}

func MakeObjekteInflator[
	T schnittstellen.Objekte[T],
	T1 schnittstellen.ObjektePtr[T],
	T2 any,
	T3 schnittstellen.VerzeichnissePtr[T2, T],
](
	or schnittstellen.ObjekteReaderFactory,
	ar schnittstellen.AkteReaderFactory,
	objekteParser schnittstellen.Parser[T, T1],
	akteParser schnittstellen.Parser[T, T1],
	pool collections.PoolLike[T],
) *objekteInflator[T, T1, T2, T3] {
	if objekteParser == nil {
		objekteParser = MakeFormat[T, T1]()
	}

	return &objekteInflator[T, T1, T2, T3]{
		or:            or,
		ar:            ar,
		objekteParser: objekteParser,
		akteParser:    akteParser,
		pool:          pool,
	}
}

func (h *objekteInflator[T, T1, T2, T3]) InflateObjekteFromSku(
	sk sku.Sku,
) (o T1, err error) {
	if h.pool == nil {
		o = T1(new(T))
	} else {
		o = h.pool.Get()
	}

	if err = h.readObjekte(sk, o); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh := o.GetAkteSha()

	if sh.IsNull() {
		return
	}

	if err = h.readAkte(sh, o); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (h *objekteInflator[T, T1, T2, T3]) readObjekte(
	sk sku.DataIdentity,
	o T1,
) (err error) {
	var r sha.ReadCloser

	if r, err = h.or.ObjekteReader(sk, sk.GetObjekteSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, r)

	var n int64

	if n, err = h.objekteParser.Parse(r, o); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.Log().Printf("parsed %d objekte bytes", n)

	return
}

func (h *objekteInflator[T, T1, T2, T3]) readAkte(
	sh schnittstellen.Sha,
	o T1,
) (err error) {
	if h.akteParser == nil {
		return
	}

	var r sha.ReadCloser

	if r, err = h.ar.AkteReader(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, r)

	var n int64

	if n, err = h.akteParser.Parse(r, o); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.Log().Printf("parsed %d akte bytes", n)

	return
}
