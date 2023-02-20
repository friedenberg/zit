package etikett

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
	Sha  sha.Sha
	Akte Akte
}

func (_ Objekte) GetGattung() schnittstellen.Gattung {
	return gattung.Etikett
}

func (o Objekte) GetAkteSha() schnittstellen.Sha {
	return o.Sha
}

func (o *Objekte) SetAkteSha(v schnittstellen.Sha) {
	o.Sha = sha.Make(v)
}

func (a *Objekte) Reset() {
	a.Sha = sha.Sha{}
	a.Akte = Akte{}
}

func (a *Objekte) ResetWith(b Objekte) {
	a.Sha = b.Sha
	a.Akte = b.Akte
}

func (a Objekte) Equals(b Objekte) bool {
	if !a.Sha.Equals(b.Sha) {
		return false
	}

	if a.Akte != b.Akte {
		return false
	}

	return true
}
