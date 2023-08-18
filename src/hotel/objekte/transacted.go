package objekte

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/golf/sku"
)

type Transacted[
	T Akte[T],
	T1 AktePtr[T],
	T2 kennung.KennungLike[T2],
	T3 kennung.KennungLikePtr[T2],
] struct {
	Sku sku.Transacted[T2, T3]
}

func (t Transacted[T, T1, T2, T3]) Kennung() T3 {
	return &t.Sku.Kennung
}

func (t Transacted[T, T1, T2, T3]) GetMetadatei() metadatei.Metadatei {
	return t.Sku.GetMetadatei()
}

func (t *Transacted[T, T1, T2, T3]) GetMetadateiPtr() *metadatei.Metadatei {
	return t.Sku.GetMetadateiPtr()
}

func (t *Transacted[T, T1, T2, T3]) SetMetadatei(
	m metadatei.Metadatei,
) {
	t.GetMetadateiPtr().ResetWith(m)
	t.SetAkteSha(m.AkteSha)
	t.SetTai(m.Tai)
}

func (t Transacted[T, T1, T2, T3]) GetAkteSha() schnittstellen.ShaLike {
	return t.Sku.Metadatei.AkteSha
}

func (t *Transacted[T, T1, T2, T3]) SetAkteSha(
	s schnittstellen.ShaLike,
) {
	sh := sha.Make(s)
	t.GetMetadateiPtr().AkteSha = sh
	t.Sku.Metadatei.AkteSha = sh
}

func (t Transacted[T, T1, T2, T3]) GetObjekteSha() schnittstellen.ShaLike {
	if !t.GetAkteSha().IsNull() && t.Sku.ObjekteSha.IsNull() {
		errors.Todo(
			"objekte sha is null while akte sha is %s",
			t.GetAkteSha(),
		)
	}

	return t.Sku.ObjekteSha
}

func (t *Transacted[T, T1, T2, T3]) SetObjekteSha(
	sh schnittstellen.ShaLike,
) {
	t.Sku.ObjekteSha = sha.Make(sh)
}

func (t Transacted[T, T1, T2, T3]) GetTai() kennung.Tai {
	return t.Sku.GetTai()
}

func (t *Transacted[T, T1, T2, T3]) SetTai(ta kennung.Tai) {
	t.Sku.SetTai(ta)
}

func (t Transacted[T, T1, T2, T3]) GetGattung() schnittstellen.GattungLike {
	return t.Sku.GetGattung()
}

func (zt Transacted[T, T1, T2, T3]) IsNew() bool {
	return zt.Sku.Kopf.Equals(zt.Sku.GetTai())
}

func (a Transacted[T, T1, T2, T3]) Less(
	b Transacted[T, T1, T2, T3],
) bool {
	return a.Sku.Less(b.Sku)
}

func (a Transacted[T, T1, T2, T3]) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a Transacted[T, T1, T2, T3]) Equals(
	b Transacted[T, T1, T2, T3],
) bool {
	if !a.GetMetadatei().Equals(b.GetMetadatei()) {
		return false
	}

	if !a.Sku.Equals(b.Sku) {
		return false
	}

	return true
}

func (a Transacted[T, T1, T2, T3]) GetEtiketten() kennung.EtikettSet {
	return a.GetMetadatei().GetEtiketten()
}

func (a Transacted[T, T1, T2, T3]) GetTyp() (t kennung.Typ) {
	tgs := []any{
		// a.Verzeichnisse,
		a.GetMetadatei(),
	}

	for _, o := range tgs {
		if tg, ok := o.(kennung.TypGetter); ok {
			t = tg.GetTyp()
			return
		}
	}

	if kennung.IsEmpty(t) {
		panic("typ is empty")
	}

	return
}

// func (a Transacted[T, T1, T2, T3]) GetIdLike() (il kennung.Kennung) {
// 	return a.Sku.GetKennung()
// }

func (a Transacted[T, T1, T2, T3]) GetSkuLike() (sk sku.SkuLike) {
	return a.Sku
}

func (a *Transacted[T, T1, T2, T3]) GetSkuLikePtr() (sk sku.SkuLikePtr) {
	return &a.Sku
}

func (a Transacted[T, T1, T2, T3]) String() string {
	return a.GetSkuLike().String()
}

func (a Transacted[T, T1, T2, T3]) GetKennungLike() kennung.Kennung {
	return a.Sku.GetKennung()
}

func (a *Transacted[T, T1, T2, T3]) GetKennungPtr() kennung.KennungPtr {
	return T3(&a.Sku.Kennung)
}

func (a *Transacted[T, T1, T2, T3]) SetKennungLike(
	k1 kennung.Kennung,
) (err error) {
	var k T2

	if err = T3(&k).Set(k1.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	a.Sku.Kennung = k

	return
}

func (a *Transacted[T, T1, T2, T3]) SetSkuLike(
	o sku.SkuLike,
) (err error) {
	return a.Sku.SetFromSkuLike(o)
	// var h T2
	// ok := false

	// if h, ok = o.GetId().(T2); !ok {
	// 	err = errors.Errorf(
	// 		"wrong type for Kennung. Expected %T but got %T",
	// 		h,
	// 		o.GetId(),
	// 	)
	// 	return
	// }

	// a.Sku.Kennung = h
	// a.Sku.ObjekteSha = sha.Make(o.GetObjekteSha())
	// // a.Sku.TransactionIndex = o.GetTransactionIndex()
	// // TODO-P3 fix sku kopf and schwanz
	// // a.Sku.Kopf = t
	// a.SetTai(o.GetTai())

	// return
}

func (a Transacted[T, T1, T2, T3]) GetKennungString() string {
	return a.Sku.GetKennung().String()
}

func (a *Transacted[T, T1, T2, T3]) Reset() {
	a.Sku.Reset()
}

func (a *Transacted[T, T1, T2, T3]) ResetWithPtr(
	b *Transacted[T, T1, T2, T3],
) {
	a.Sku.ResetWith(b.Sku)
}

func (a *Transacted[T, T1, T2, T3]) ResetWith(
	b Transacted[T, T1, T2, T3],
) {
	a.Sku.ResetWith(b.Sku)
}
