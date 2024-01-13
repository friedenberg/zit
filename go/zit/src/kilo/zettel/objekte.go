package zettel

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/gattung"
)

type Objekte struct{}

func (z Objekte) GetGattung() schnittstellen.GattungLike {
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
