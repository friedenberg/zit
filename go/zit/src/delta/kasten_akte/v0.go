package kasten_akte

import (
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/values"
	"code.linenisgreat.com/zit/src/charlie/gattung"
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
