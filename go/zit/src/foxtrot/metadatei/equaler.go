package metadatei

import (
	"code.linenisgreat.com/zit-go/src/bravo/iter"
	"code.linenisgreat.com/zit-go/src/echo/kennung"
)

var (
	EqualerSansTai equalerSansTai
	Equaler        equaler
)

type equalerSansTai struct{}

func (equalerSansTai) Equals(pz, z1 *Metadatei) bool {
	if !pz.Akte.EqualsSha(&z1.Akte) {
		return false
	}

	if !pz.Typ.Equals(z1.Typ) {
		return false
	}

	if !iter.SetEquals[kennung.Etikett](pz.GetEtiketten(), z1.GetEtiketten()) {
		return false
	}

	if !pz.Bezeichnung.Equals(z1.Bezeichnung) {
		return false
	}

	return true
}

type equaler struct{}

func (equaler) Equals(pz, z1 *Metadatei) bool {
	if !EqualerSansTai.Equals(pz, z1) {
		return false
	}

	if !pz.Tai.Equals(z1.Tai) {
		return false
	}

	return true
}
