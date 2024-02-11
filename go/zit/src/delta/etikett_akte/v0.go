package etikett_akte

import (
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/charlie/gattung"
)

type V0 struct{}

func (V0) GetGattung() schnittstellen.GattungLike {
	return gattung.Etikett
}

func (ct *V0) Merge(ct2 V0) {
}

func (a *V0) Reset() {
}

func (a *V0) ResetWith(b V0) {
}

func (a V0) Equals(b V0) bool {
	return true
}
