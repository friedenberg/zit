package metadatei

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type WithKennungLike interface {
	GetKennung() kennung.IdLike
	Getter
}

type WithKennungPtrLike interface {
	WithKennungLike
	GetterPtr
	Setter
}

type WithKennung struct {
	Kennung   kennung.KennungPtr
	Metadatei Metadatei
}

func (a WithKennung) String() string {
	return a.Kennung.String()
}

func (a WithKennung) GetKennung() kennung.IdLike {
	return a.Kennung
}

func (a WithKennung) GetGattung() schnittstellen.Gattung {
	return a.Kennung.GetGattung()
}

func (a WithKennung) GetMetadatei() Metadatei {
	return a.Metadatei
}

func (a *WithKennung) GetMetadateiPtr() *Metadatei {
	return &a.Metadatei
}

func (a *WithKennung) SetMetadatei(m Metadatei) {
	a.Metadatei.SetMetadatei(m)
}

func (a WithKennung) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a WithKennung) Equals(b WithKennung) bool {
	if a.Kennung.String() != b.Kennung.String() {
		return false
	}

	if !a.Metadatei.Equals(b.Metadatei) {
		return false
	}

	return true
}

func (wk *WithKennung) Reset() {
	wk.Kennung.Reset()
	wk.Metadatei.Reset()
}
