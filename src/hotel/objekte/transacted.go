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
	T4 gattung.Verzeichnisse[T],
	T5 gattung.VerzeichnissePtr[T4, T],
] struct {
	Objekte       T
	Verzeichnisse T4
	Sku           sku.Transacted[T2, T3]
}

func (t Transacted[T, T1, T2, T3, T4, T5]) Stored() *Stored[T, T1, T2, T3] {
	return &Stored[T, T1, T2, T3]{
		Objekte: t.Objekte,
		Sku: sku.External[T2, T3]{
			ObjekteSha: t.Sku.ObjekteSha,
			AkteSha:    t.Sku.AkteSha,
			Kennung:    t.Sku.Kennung,
		},
	}
}

func (t Transacted[T, T1, T2, T3, T4, T5]) Kennung() T3 {
	return &t.Sku.Kennung
}

func (t Transacted[T, T1, T2, T3, T4, T5]) AkteSha() sha.Sha {
	return t.Objekte.GetAkteSha()
}

func (t *Transacted[T, T1, T2, T3, T4, T5]) SetAkteSha(
	s sha.Sha,
) {
	T1(&t.Objekte).SetAkteSha(s)
}

func (t Transacted[T, T1, T2, T3, T4, T5]) ObjekteSha() sha.Sha {
	return t.Sku.ObjekteSha
}

func (t *Transacted[T, T1, T2, T3, T4, T5]) SetObjekteSha(
	arf gattung.AkteReaderFactory,
	v string,
) (err error) {
	if err = t.Sku.ObjekteSha.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	//TODO-P3 provide opportunity for Objekte to bootstrap from arf

	return
}

func (t Transacted[T, T1, T2, T3, T4, T5]) GetGattung() gattung.Gattung {
	return t.Sku.Kennung.GetGattung()
}

func (zt Transacted[T, T1, T2, T3, T4, T5]) IsNew() bool {
	return zt.Sku.Kopf == zt.Sku.Schwanz && zt.Sku.TransactionIndex == 0
}

func (a Transacted[T, T1, T2, T3, T4, T5]) Less(
	b Transacted[T, T1, T2, T3, T4, T5],
) bool {
	return a.Sku.GetTime().Less(b.Sku.GetTime())
}

func (a Transacted[T, T1, T2, T3, T4, T5]) Equals(
	b *Transacted[T, T1, T2, T3, T4, T5],
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

func (a *Transacted[T, T1, T2, T3, T4, T5]) SetSkuLike(
	o sku.SkuLike,
) (err error) {
	var h T2
	ok := false

	if h, ok = o.GetId().(T2); !ok {
		err = errors.Errorf("wrong type for Kennung. Expected %T but got %T", h, o.GetId())
		return
	}

	a.Sku.Kennung = h
	a.Sku.ObjekteSha = o.GetObjekteSha()
	a.Sku.TransactionIndex = o.GetTransactionIndex()
	//TODO-P3 fix sku kopf and schwanz
	// a.Sku.Kopf = t
	a.Sku.Schwanz = o.GetTime()
	a.Sku.Mutter = o.GetMutter()

	return
}

func (a Transacted[T, T1, T2, T3, T4, T5]) GetKennungString() string {
	return a.Sku.Kennung.String()
}

func (a *Transacted[T, T1, T2, T3, T4, T5]) GenerateVerzeichnisse() {
	T5(&a.Verzeichnisse).ResetWithObjekte(&a.Objekte)
}

func (a *Transacted[T, T1, T2, T3, T4, T5]) Reset(
	b *Transacted[T, T1, T2, T3, T4, T5],
) {
	if b == nil {
		a.Sku.Reset(nil)
		T1(&a.Objekte).Reset(nil)
	} else {
		a.Sku.Reset(&b.Sku)
		T1(&a.Objekte).Reset(&b.Objekte)
	}

	a.GenerateVerzeichnisse()
}
