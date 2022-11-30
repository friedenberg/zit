package objekte

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/echo/sku"
)

type Transacted[
	T gattung.Objekte2,
	T1 gattung.Objekte2Ptr[T],
	T2 gattung.Identifier2[T2],
	T3 gattung.IdentifierPtr[T2],
] struct {
	Objekte T
	Sku     sku.Sku2[T2, T3]
}

func (t Transacted[T, T1, T2, T3]) Kennung() T3 {
	return &t.Sku.Kennung
}

func (t Transacted[T, T1, T2, T3]) AkteSha() sha.Sha {
	return t.Objekte.AkteSha()
}

func (t *Transacted[T, T1, T2, T3]) SetAkteSha(
	s sha.Sha,
) {
	T1(&t.Objekte).SetAkteSha(s)
}

func (t Transacted[T, T1, T2, T3]) ObjekteSha() sha.Sha {
	return t.Sku.Sha
}

func (t *Transacted[T, T1, T2, T3]) SetObjekteSha(
	arf gattung.AkteReaderFactory,
	v string,
) (err error) {
	if err = t.Sku.Sha.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	//TODO-P0 provide opportunity for Objekte to bootstrap from arf

	return
}

func (t Transacted[T, T1, T2, T3]) Gattung() gattung.Gattung {
	return t.Sku.Kennung.Gattung()
}

func (zt Transacted[T, T1, T2, T3]) IsNew() bool {
	return zt.Sku.Kopf == zt.Sku.Schwanz && zt.Sku.TransactionIndex == 0
}

// TODO-P0
func (a Transacted[T, T1, T2, T3]) Equals(
	b Transacted[T, T1, T2, T3],
) bool {
	return false
}

func (a *Transacted[T, T1, T2, T3]) Reset(
	b *Transacted[T, T1, T2, T3],
) {
	if b == nil {
		a.Sku.Reset(nil)
		T1(&a.Objekte).Reset(nil)
	} else {
		a.Sku.Reset(&b.Sku)
		T1(&a.Objekte).Reset(&b.Objekte)
	}
}
