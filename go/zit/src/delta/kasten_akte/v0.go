package kasten_akte

import (
	"code.linenisgreat.com/zit-go/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit-go/src/bravo/values"
	"code.linenisgreat.com/zit-go/src/charlie/gattung"
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
