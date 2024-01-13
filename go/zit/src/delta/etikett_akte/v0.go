package etikett_akte

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/gattung"
)

type V0 struct{}

func (_ V0) GetGattung() schnittstellen.GattungLike {
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
