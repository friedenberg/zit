package erworben

import (
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/objekte"
	"github.com/friedenberg/zit/src/schnittstellen"
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
	//TODO-P3 rename to AkteSha
	Sha  sha.Sha
	Akte Akte
}

func (o Objekte) GetAkteSha() schnittstellen.Sha {
	return o.Sha
}

func (o *Objekte) SetAkteSha(v schnittstellen.Sha) {
	o.Sha = sha.Make(v)
}

func (a Objekte) Equals(b *Objekte) bool {
	if b == nil {
		return false
	}

	if !a.Sha.Equals(b.Sha) {
		return false
	}

	return true
}

func (a *Objekte) Reset(b *Objekte) {
	if b == nil {
		a.Sha = sha.Sha{}
		a.Akte.Reset(nil)
	} else {
		a.Sha = b.Sha
		a.Akte.Reset(&b.Akte)
	}
}

func (c Objekte) GetGattung() schnittstellen.Gattung {
	return gattung.Konfig
}
