package objekte

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/golf/sku"
)

type Transacted[
	T gattung.Objekte[T],
	T1 gattung.ObjektePtr[T],
	T2 gattung.Identifier[T2],
	T3 gattung.IdentifierPtr[T2],
] struct {
	Objekte T
	Sku     sku.Transacted[T2, T3]
}

func (t Transacted[T, T1, T2, T3]) Stored() *Stored[T, T1, T2, T3] {
	return &Stored[T, T1, T2, T3]{
		Objekte: t.Objekte,
		Sku: sku.External[T2, T3]{
			Sha:     t.Sku.Sha,
			Kennung: t.Sku.Kennung,
		},
	}
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

	//TODO-P3 provide opportunity for Objekte to bootstrap from arf

	return
}

func (t Transacted[T, T1, T2, T3]) Gattung() gattung.Gattung {
	return t.Sku.Kennung.Gattung()
}

func (zt Transacted[T, T1, T2, T3]) IsNew() bool {
	return zt.Sku.Kopf == zt.Sku.Schwanz && zt.Sku.TransactionIndex == 0
}

func (a Transacted[T, T1, T2, T3]) Equals(
	b *Transacted[T, T1, T2, T3],
) bool {
	if b == nil {
		return false
	}

	if !a.Sku.Equals(&b.Sku) {
		return false
	}

	if !a.Objekte.Equals(&b.Objekte) {
		return false
	}

	return true
}

func (a *Transacted[T, T1, T2, T3]) SetSku(
	v sku.Transacted[T2, T3],
) {
	a.Sku = v
}

func (a Transacted[T, T1, T2, T3]) GetSku() (v sku.Transacted[T2, T3]) {
	v = a.Sku
  return
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