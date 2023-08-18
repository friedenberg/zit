package objekte

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/checkout_mode"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/golf/sku"
)

type ExternalKeyer[
	T Akte[T],
	T1 AktePtr[T],
	T2 kennung.KennungLike[T2],
	T3 kennung.KennungLikePtr[T2],
] struct{}

func (_ ExternalKeyer[T, T1, T2, T3]) Key(e *External[T, T1, T2, T3]) string {
	if e == nil {
		return ""
	}

	return e.Sku.GetKennung().String()
}

type External[
	T Akte[T],
	T1 AktePtr[T],
	T2 kennung.KennungLike[T2],
	T3 kennung.KennungLikePtr[T2],
] struct {
	Sku  sku.External[T2, T3]
}

func (a External[T, T1, T2, T3]) GetSkuLike() sku.SkuLike {
	return a.Sku
}

func (a *External[T, T1, T2, T3]) GetSkuLikePtr() sku.SkuLikePtr {
	return &a.Sku
}

func (a External[T, T1, T2, T3]) GetMetadatei() metadatei.Metadatei {
	return a.Sku.GetMetadatei()
}

func (a *External[T, T1, T2, T3]) GetMetadateiPtr() *metadatei.Metadatei {
	return a.Sku.GetMetadateiPtr()
}

func (a *External[T, T1, T2, T3]) SetMetadatei(m metadatei.Metadatei) {
	a.GetMetadateiPtr().ResetWith(m)
	a.SetAkteSha(m.AkteSha)
}

func (a External[T, T1, T2, T3]) GetEtiketten() kennung.EtikettSet {
	egs := []any{
		// a.Verzeichnisse,
		a.GetMetadatei(),
	}

	for _, o := range egs {
		if eg, ok := o.(kennung.EtikettenGetter); ok {
			return eg.GetEtiketten()
		}
	}

	return kennung.MakeEtikettSet()
}

func (a External[T, T1, T2, T3]) GetTyp() (t kennung.Typ) {
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

	return
}

func (a External[T, T1, T2, T3]) GetKennungLike() kennung.Kennung {
	return a.Sku.GetKennung()
}

func (a *External[T, T1, T2, T3]) GetKennungPtr() kennung.KennungPtr {
	return T3(&a.Sku.Kennung)
}

func (a External[T, T1, T2, T3]) GetIdLike() (il kennung.Kennung) {
	return a.Sku.GetKennung()
}

func (a External[T, T1, T2, T3]) String() string {
	return a.Sku.String()
}

func (a External[T, T1, T2, T3]) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a External[T, T1, T2, T3]) Equals(b External[T, T1, T2, T3]) bool {
	if !a.Sku.Equals(b.Sku) {
		return false
	}

	return true
}

func (e External[T, T1, T2, T3]) GetGattung() schnittstellen.GattungLike {
	return e.Sku.GetGattung()
}

func (e External[T, T1, T2, T3]) GetFDs() sku.ExternalFDs {
	return e.Sku.FDs
}

func (e External[T, T1, T2, T3]) GetObjekteFD() kennung.FD {
	return e.Sku.FDs.Objekte
}

func (e External[T, T1, T2, T3]) GetAkteFD() kennung.FD {
	return e.Sku.FDs.Akte
}

func (e *External[T, T1, T2, T3]) GetFDsPtr() *sku.ExternalFDs {
	return &e.Sku.FDs
}

func (e External[T, T1, T2, T3]) GetAktePath() string {
	return e.Sku.FDs.Akte.Path
}

func (e External[T, T1, T2, T3]) GetObjekteSha() schnittstellen.ShaLike {
	return e.Sku.ObjekteSha
}

func (t External[T, T1, T2, T3]) GetFDAkteSha() schnittstellen.ShaLike {
	return t.Sku.FDs.Akte.Sha
}

func (t External[T, T1, T2, T3]) GetAkteSha() schnittstellen.ShaLike {
	return t.Sku.GetAkteSha()
}

func (e *External[T, T1, T2, T3]) SetAkteFD(fd kennung.FD) {
	e.Sku.FDs.Akte = fd
	e.SetAkteSha(fd.Sha)
}

func (e *External[T, T1, T2, T3]) SetAkteSha(v schnittstellen.ShaLike) {
	sh := sha.Make(v)
	e.GetMetadateiPtr().AkteSha = sh
	e.Sku.SetAkteSha(sh)
}

func (e *External[T, T1, T2, T3]) SetObjekteSha(
	sh schnittstellen.ShaLike,
) {
	e.Sku.ObjekteSha = sha.Make(sh)
}

func (e External[T, T1, T2, T3]) GetCheckoutMode() (checkout_mode.Mode, error) {
	return e.Sku.GetCheckoutMode()
}
