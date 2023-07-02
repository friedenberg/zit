package sku

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

type WithKennungLike interface {
	GetKennungLike() kennung.Kennung
	metadatei.Getter
}

type WithKennungPtrLike interface {
	WithKennungLike
	metadatei.GetterPtr
	metadatei.Setter
}

type WithKennungInterface struct {
	Kennung   kennung.Kennung
	Metadatei Metadatei
}

func (a WithKennungInterface) String() string {
	return a.Kennung.String()
}

func (a WithKennungInterface) GetKennungLike() kennung.Kennung {
	return a.Kennung
}

func (a WithKennungInterface) GetGattung() schnittstellen.GattungLike {
	return a.Kennung.GetGattung()
}

func (a WithKennungInterface) GetMetadatei() Metadatei {
	return a.Metadatei
}

func (a *WithKennungInterface) GetMetadateiPtr() *Metadatei {
	return &a.Metadatei
}

func (a *WithKennungInterface) SetMetadatei(m Metadatei) {
	a.Metadatei.SetMetadatei(m)
}

func (a WithKennungInterface) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a WithKennungInterface) Equals(b WithKennungInterface) bool {
	if a.Kennung.String() != b.Kennung.String() {
		return false
	}

	if !a.Metadatei.Equals(b.Metadatei) {
		return false
	}

	return true
}

func (wk *WithKennungInterface) Reset() {
	k := wk.Kennung.KennungPtrClone()
	k.Reset()
	wk.Kennung = k.KennungClone()
	wk.Metadatei.Reset()
}
