package metadatei

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type WithKennung[T kennung.KennungLike[T], T1 kennung.KennungLikePtr[T]] struct {
	Kennung   T
	Metadatei Metadatei
}

func (a WithKennung[K, KPtr]) String() string {
	return a.Kennung.String()
}

func (a WithKennung[K, KPtr]) GetKennungLike() kennung.Kennung {
	return a.Kennung
}

func (a WithKennung[K, KPtr]) GetKennung() K {
	return a.Kennung
}

func (a WithKennung[K, KPtr]) GetGattung() schnittstellen.Gattung {
	return a.Kennung.GetGattung()
}

func (a WithKennung[K, KPtr]) GetMetadatei() Metadatei {
	return a.Metadatei
}

func (a *WithKennung[K, KPtr]) GetMetadateiPtr() *Metadatei {
	return &a.Metadatei
}

func (a *WithKennung[K, KPtr]) SetMetadatei(m Metadatei) {
	a.Metadatei.SetMetadatei(m)
}

func (a WithKennung[K, KPtr]) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a WithKennung[K, KPtr]) Equals(b WithKennung[K, KPtr]) bool {
	if a.Kennung.String() != b.Kennung.String() {
		return false
	}

	if !a.Metadatei.Equals(b.Metadatei) {
		return false
	}

	return true
}

func (wk *WithKennung[K, KPtr]) Reset() {
	KPtr(&wk.Kennung).Reset()
	wk.Metadatei.Reset()
}
