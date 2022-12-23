package objekte

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/golf/transaktion"
)

type Inflator[
	T gattung.Objekte[T],
	T1 gattung.ObjektePtr[T],
	T2 gattung.Identifier[T2],
	T3 gattung.IdentifierPtr[T2],
] interface {
	Inflate(*transaktion.Transaktion, *sku.Sku) (*Transacted[T, T1, T2, T3], error)
}

type transactedInflator[
	T gattung.Objekte[T],
	T1 gattung.ObjektePtr[T],
	T2 gattung.Identifier[T2],
	T3 gattung.IdentifierPtr[T2],
] struct {
	arf              gattung.AkteReaderFactory
	frc              FuncReadCloser
	objekteFormatter Formatter2
	akteParser       gattung.Parser[T, T1]
}

func MakeTransactedInflator[
	T gattung.Objekte[T],
	T1 gattung.ObjektePtr[T],
	T2 gattung.Identifier[T2],
	T3 gattung.IdentifierPtr[T2],
](
	arf gattung.AkteReaderFactory,
	frc FuncReadCloser,
	akteParser gattung.Parser[T, T1],
) *transactedInflator[T, T1, T2, T3] {
	return &transactedInflator[T, T1, T2, T3]{
		arf:              arf,
		frc:              frc,
		objekteFormatter: *MakeFormatter2(),
		akteParser:       akteParser,
	}
}

func (h *transactedInflator[T, T1, T2, T3]) Inflate(
	tr *transaktion.Transaktion,
	o *sku.Sku,
) (t *Transacted[T, T1, T2, T3], err error) {
	t = new(Transacted[T, T1, T2, T3])

	if err = t.SetTransactionAndObjekte(
		tr,
		o,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	{
		var r sha.ReadCloser

		if r, err = h.frc(t.ObjekteSha()); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.Deferred(&err, r.Close)

		if _, err = h.objekteFormatter.ReadFormat(r, t); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if h.akteParser != nil {
		var r sha.ReadCloser

		if r, err = h.arf.AkteReader(t.AkteSha()); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.Deferred(&err, r.Close)

		if _, err = h.akteParser.Parse(r, &t.Objekte); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
