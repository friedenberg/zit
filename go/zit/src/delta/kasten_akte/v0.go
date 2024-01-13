package kasten_akte

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/charlie/gattung"
)

type V0 struct {
	Uri values.Uri `toml:"uri"`
}

func (_ V0) GetGattung() schnittstellen.GattungLike {
	return gattung.Typ
}

func (a *V0) Reset() {
	a.Uri = values.Uri{}
}

func (a *V0) ResetWith(b V0) {
	a.Uri = b.Uri
}

func (a V0) Equals(b V0) bool {
	if a.Uri != b.Uri {
		return false
	}

	return true
}
