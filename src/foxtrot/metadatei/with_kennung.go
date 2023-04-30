package metadatei

import "github.com/friedenberg/zit/src/delta/kennung"

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

func (a WithKennung) GetKennung() kennung.IdLike {
	return a.Kennung
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
