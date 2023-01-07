package objekte

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/foxtrot/ts"
	"github.com/friedenberg/zit/src/golf/sku"
)

type TransactedInflator[
	T gattung.Objekte[T],
	T1 gattung.ObjektePtr[T],
	T2 gattung.Identifier[T2],
	T3 gattung.IdentifierPtr[T2],
	T4 gattung.Verzeichnisse[T],
	T5 gattung.VerzeichnissePtr[T4, T],
] interface {
	Inflate(ts.Time, sku.SkuLike) (*Transacted[T, T1, T2, T3, T4, T5], error)
}

type transactedInflator[
	T gattung.Objekte[T],
	T1 gattung.ObjektePtr[T],
	T2 gattung.Identifier[T2],
	T3 gattung.IdentifierPtr[T2],
	T4 gattung.Verzeichnisse[T],
	T5 gattung.VerzeichnissePtr[T4, T],
] struct {
	orc           sku.FuncSkuObjekteReader
	arc           gattung.FuncReadCloser
	objekteParser gattung.Parser[T, T1]
	akteParser    gattung.Parser[T, T1]
	pool          collections.PoolLike[Transacted[T, T1, T2, T3, T4, T5]]
}

func MakeTransactedInflator[
	T gattung.Objekte[T],
	T1 gattung.ObjektePtr[T],
	T2 gattung.Identifier[T2],
	T3 gattung.IdentifierPtr[T2],
	T4 gattung.Verzeichnisse[T],
	T5 gattung.VerzeichnissePtr[T4, T],
](
	orc sku.FuncSkuObjekteReader,
	arc gattung.FuncReadCloser,
	objekteParser gattung.Parser[T, T1],
	akteParser gattung.Parser[T, T1],
	pool collections.PoolLike[Transacted[T, T1, T2, T3, T4, T5]],
) *transactedInflator[T, T1, T2, T3, T4, T5] {
	if objekteParser == nil {
		objekteParser = MakeFormat[T, T1]()
	}

	return &transactedInflator[T, T1, T2, T3, T4, T5]{
		orc:           orc,
		arc:           arc,
		objekteParser: objekteParser,
		akteParser:    akteParser,
		pool:          pool,
	}
}

// TODO-P3 rename to InflateFromSku
func (h *transactedInflator[T, T1, T2, T3, T4, T5]) Inflate(
	ti ts.Time,
	o sku.SkuLike,
) (t *Transacted[T, T1, T2, T3, T4, T5], err error) {
	if h.pool == nil {
		t = new(Transacted[T, T1, T2, T3, T4, T5])
	} else {
		t = h.pool.Get()
	}

	if err = t.SetTimeAndObjekte(ti, o); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = h.readObjekte(o, t); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh := t.AkteSha()

	if sh.IsNull() {
		return
	}

	if err = h.readAkte(sh, t); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (h *transactedInflator[T, T1, T2, T3, T4, T5]) readObjekte(
	sk sku.SkuLike,
	t *Transacted[T, T1, T2, T3, T4, T5],
) (err error) {
	var r sha.ReadCloser

	if r, err = h.orc(sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, r.Close)

	var n int64

	if n, err = h.objekteParser.Parse(r, &t.Objekte); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.Log().Printf("parsed %d objekte bytes", n)

	return
}

func (h *transactedInflator[T, T1, T2, T3, T4, T5]) readAkte(
	sh sha.Sha,
	t *Transacted[T, T1, T2, T3, T4, T5],
) (err error) {
	if h.akteParser == nil {
		return
	}

	var r sha.ReadCloser

	if r, err = h.arc(t.AkteSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, r)

	var n int64

	if n, err = h.akteParser.Parse(r, &t.Objekte); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.Log().Printf("parsed %d akte bytes", n)

	return
}
