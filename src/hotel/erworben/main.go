package erworben

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
	kennung.Konfig,
	*kennung.Konfig,
	objekte.NilVerzeichnisse[Objekte],
	*objekte.NilVerzeichnisse[Objekte],
]

type Objekte struct {
	Akte
}

func (a Objekte) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a Objekte) Equals(b Objekte) bool {
	return true
}

func (a *Objekte) Reset() {
	a.Akte.Reset()
}

func (a *Objekte) ResetWith(b Objekte) {
	a.Akte.ResetWith(b.Akte)
}

func (c Objekte) GetGattung() schnittstellen.Gattung {
	return gattung.Konfig
}
