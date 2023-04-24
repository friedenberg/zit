package etikett

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/objekte"
)

type Transacted = objekte.Transacted[
	Objekte,
	*Objekte,
	kennung.Etikett,
	*kennung.Etikett,
	objekte.NilVerzeichnisse[Objekte],
	*objekte.NilVerzeichnisse[Objekte],
]

type ExternalKeyer = objekte.ExternalKeyer[
	Objekte,
	*Objekte,
	kennung.Etikett,
	*kennung.Etikett,
]

type External = objekte.External[
	Objekte,
	*Objekte,
	kennung.Etikett,
	*kennung.Etikett,
]

type CheckedOut = objekte.CheckedOut[
	Objekte,
	*Objekte,
	kennung.Etikett,
	*kennung.Etikett,
	objekte.NilVerzeichnisse[Objekte],
	*objekte.NilVerzeichnisse[Objekte],
]

type Objekte struct {
	Akte
}

func (a Objekte) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (_ Objekte) GetGattung() schnittstellen.Gattung {
	return gattung.Etikett
}

func (a *Objekte) Reset() {
	a.Akte = Akte{}
}

func (a *Objekte) ResetWith(b Objekte) {
	a.Akte = b.Akte
}

func (a Objekte) Equals(b Objekte) bool {
	if a.Akte != b.Akte {
		return false
	}

	return true
}
