package kasten

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/values"
)

type Akte struct {
	Uri values.Uri `toml:"uri"`
}

func (_ Akte) GetGattung() schnittstellen.Gattung {
	return gattung.Typ
}

func (a *Akte) Reset() {
	a.Uri = values.Uri{}
}

func (a *Akte) ResetWith(b Akte) {
	a.Uri = b.Uri
}

func (a Akte) Equals(b Akte) bool {
	if a.Uri != b.Uri {
		return false
	}

	return true
}
