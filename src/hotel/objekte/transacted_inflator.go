package objekte

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
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
	Inflate(ts.Time, *sku.Sku) (*Transacted[T, T1, T2, T3, T4, T5], error)
}

type transactedInflator[
	T gattung.Objekte[T],
	T1 gattung.ObjektePtr[T],
	T2 gattung.Identifier[T2],
	T3 gattung.IdentifierPtr[T2],
	T4 gattung.Verzeichnisse[T],
	T5 gattung.VerzeichnissePtr[T4, T],
] struct {
	arf           gattung.AkteReaderFactory
	frc           gattung.FuncReadCloser
	objekteParser gattung.Parser[T, T1]
	akteParser    gattung.Parser[T, T1]
}

func MakeTransactedInflator[
	T gattung.Objekte[T],
	T1 gattung.ObjektePtr[T],
	T2 gattung.Identifier[T2],
	T3 gattung.IdentifierPtr[T2],
	T4 gattung.Verzeichnisse[T],
	T5 gattung.VerzeichnissePtr[T4, T],
](
	arf gattung.AkteReaderFactory,
	frc gattung.FuncReadCloser,
	objekteParser gattung.Parser[T, T1],
	akteParser gattung.Parser[T, T1],
) *transactedInflator[T, T1, T2, T3, T4, T5] {
	if objekteParser == nil {
		objekteParser = MakeFormat[T, T1]()
	}

	return &transactedInflator[T, T1, T2, T3, T4, T5]{
		arf:           arf,
		frc:           frc,
		objekteParser: objekteParser,
		akteParser:    akteParser,
	}
}

func (h *transactedInflator[T, T1, T2, T3, T4, T5]) Inflate(
	ti ts.Time,
	o *sku.Sku,
) (t *Transacted[T, T1, T2, T3, T4, T5], err error) {
	t = new(Transacted[T, T1, T2, T3, T4, T5])

	if err = t.SetTimeAndObjekte(ti, o); err != nil {
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

		if _, err = h.objekteParser.Parse(r, &t.Objekte); err != nil {
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
