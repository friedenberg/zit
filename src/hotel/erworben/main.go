package erworben

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
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
	// TODO-P3 rename to AkteSha
	Sha  sha.Sha
	Akte Akte
}

func (o Objekte) GetAkteSha() schnittstellen.Sha {
	return o.Sha
}

func (o *Objekte) SetAkteSha(v schnittstellen.Sha) {
	o.Sha = sha.Make(v)
}

func (a Objekte) Equals(b Objekte) bool {
	if !a.Sha.Equals(b.Sha) {
		return false
	}

	return true
}

func (a *Objekte) Reset() {
	a.Sha = sha.Sha{}
	a.Akte.Reset()
}

func (a *Objekte) ResetWith(b Objekte) {
	a.Sha = b.Sha
	a.Akte.ResetWith(b.Akte)
}

func (c Objekte) GetGattung() schnittstellen.Gattung {
	return gattung.Konfig
}
