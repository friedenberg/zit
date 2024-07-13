package zettel

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
)

type Objekte struct{}

func (z Objekte) GetGattung() interfaces.Genre {
	return gattung.Zettel
}

func (z Objekte) Equals(z1 Objekte) bool {
	return true
}

func (z Objekte) IsEmpty() bool {
	return true
}

func (z *Objekte) Reset() {
}

func (z *Objekte) ResetWith(z1 Objekte) {
}
