package kasten

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/values"
)

type Objekte struct {
	Akte Akte
}

func (a Objekte) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (o Objekte) Reset() {
	o.Akte.Reset()
}

func (o Objekte) ResetWith(o1 Objekte) {
	o.Akte.ResetWith(o1.Akte)
}

func (o Objekte) Equals(o1 Objekte) bool {
	if !o.Akte.Equals(o1.Akte) {
		return false
	}

	return true
}

func (o Objekte) GetGattung() schnittstellen.Gattung {
	return gattung.Typ
}
