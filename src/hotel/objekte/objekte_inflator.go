package objekte

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/golf/sku"
)

// TODO-P1 split into ObjekteInflator
type ObjekteInflator[
	T gattung.Objekte[T],
	T1 gattung.ObjektePtr[T],
	T4 gattung.Verzeichnisse[T],
	T5 gattung.VerzeichnissePtr[T4, T],
] interface {
	InflateObjekteFromSku(sku.Sku) (*T1, error)
}

type objekteInflator[
	T gattung.Objekte[T],
	T1 gattung.ObjektePtr[T],
	T2 gattung.Verzeichnisse[T],
	T3 gattung.VerzeichnissePtr[T2, T],
] struct {
	orc           gattung.FuncObjekteReader
	arc           gattung.FuncReadCloser
	objekteParser gattung.Parser[T, T1]
	akteParser    gattung.Parser[T, T1]
	pool          collections.PoolLike[T]
}

func MakeObjekteInflator[
	T gattung.Objekte[T],
	T1 gattung.ObjektePtr[T],
	T2 gattung.Verzeichnisse[T],
	T3 gattung.VerzeichnissePtr[T2, T],
](
	orc gattung.FuncObjekteReader,
	arc gattung.FuncReadCloser,
	objekteParser gattung.Parser[T, T1],
	akteParser gattung.Parser[T, T1],
	pool collections.PoolLike[T],
) *objekteInflator[T, T1, T2, T3] {
	if objekteParser == nil {
		objekteParser = MakeFormat[T, T1]()
	}

	return &objekteInflator[T, T1, T2, T3]{
		orc:           orc,
		arc:           arc,
		objekteParser: objekteParser,
		akteParser:    akteParser,
		pool:          pool,
	}
}

// TODO-P3 rename to InflateFromSku
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

	if r, err = h.orc(sk, sk.GetObjekteSha()); err != nil {
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
	sh sha.Sha,
	o T1,
) (err error) {
	if h.akteParser == nil {
		return
	}

	var r sha.ReadCloser

	if r, err = h.arc(sh); err != nil {
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
