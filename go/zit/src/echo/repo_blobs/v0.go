package repo_blobs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
)

type V0 struct {
	Uri values.Uri `toml:"uri"`
}

func (_ V0) GetGattung() interfaces.Genre {
	return genres.Type
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
