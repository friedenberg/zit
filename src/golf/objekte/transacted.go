package objekte

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/foxtrot/sku"
)

type Transacted[
	T schnittstellen.Objekte[T],
	T1 schnittstellen.ObjektePtr[T],
	T2 schnittstellen.Id[T2],
	T3 schnittstellen.IdPtr[T2],
	T4 any,
	T5 schnittstellen.VerzeichnissePtr[T4, T],
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

func (t Transacted[T, T1, T2, T3, T4, T5]) GetAkteSha() schnittstellen.Sha {
	return t.Objekte.GetAkteSha()
}

func (t *Transacted[T, T1, T2, T3, T4, T5]) SetAkteSha(
	s schnittstellen.Sha,
) {
	T1(&t.Objekte).SetAkteSha(s)
}

func (t Transacted[T, T1, T2, T3, T4, T5]) GetObjekteSha() schnittstellen.Sha {
	if !t.GetAkteSha().IsNull() && t.Sku.ObjekteSha.IsNull() {
		errors.Todo(
			"objekte sha is null while akte sha is %s",
			t.GetAkteSha(),
		)
	}

	return t.Sku.ObjekteSha
}

func (t *Transacted[T, T1, T2, T3, T4, T5]) SetObjekteSha(
	arf schnittstellen.AkteReaderFactory,
	v string,
) (err error) {
	if err = t.Sku.ObjekteSha.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	//TODO-P3 provide opportunity for Objekte to bootstrap from arf

	return
}

func (t Transacted[T, T1, T2, T3, T4, T5]) GetGattung() schnittstellen.Gattung {
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

	if !a.Objekte.Equals(b.Objekte) {
		return false
	}

	return true
}

func (a Transacted[T, T1, T2, T3, T4, T5]) GetObjekte() (o T) {
	o = a.Objekte
	return
}

func (a Transacted[T, T1, T2, T3, T4, T5]) GetSku2() (sk sku.Sku2) {
	sk = a.Sku.Sku2()
	errors.TodoP2("make certain akte sha is in sku")
	sk.AkteSha = sha.Make(a.GetAkteSha())
	return
}

func (a *Transacted[T, T1, T2, T3, T4, T5]) SetDataIdentity(
	o sku.DataIdentity,
) (err error) {
	var h T2

	if err = T3(&h).Set(o.GetId().String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	a.Sku.Kennung = h
	a.Sku.ObjekteSha = sha.Make(o.GetObjekteSha())
	a.Sku.AkteSha = sha.Make(o.GetAkteSha())
	a.Sku.Schwanz = o.GetTime()

	return
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
	a.Sku.ObjekteSha = sha.Make(o.GetObjekteSha())
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
	T5(&a.Verzeichnisse).ResetWithObjekte(a.Objekte)
}

func (a *Transacted[T, T1, T2, T3, T4, T5]) Reset() {
	a.Sku.Reset()
	T1(&a.Objekte).Reset()
	a.GenerateVerzeichnisse()
}

func (a *Transacted[T, T1, T2, T3, T4, T5]) ResetWith(
	b Transacted[T, T1, T2, T3, T4, T5],
) {
	a.Sku.ResetWith(b.Sku)
	T1(&a.Objekte).ResetWith(b.Objekte)
	a.GenerateVerzeichnisse()
}