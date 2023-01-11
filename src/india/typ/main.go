package typ

import (
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/foxtrot/kennung"
	"github.com/friedenberg/zit/src/hotel/objekte"
)

type Transacted = objekte.Transacted[
	Objekte,
	*Objekte,
	kennung.Typ,
	*kennung.Typ,
	objekte.NilVerzeichnisse[Objekte],
	*objekte.NilVerzeichnisse[Objekte],
]

type Objekte struct {
	Sha  sha.Sha
	Akte Akte
}

func (o Objekte) Reset(o1 *Objekte) {
	if o1 == nil {
		o.Sha = sha.Sha{}
		o.Akte.Reset(nil)
	} else {
		o.Sha = o1.Sha
		o.Akte.Reset(&o1.Akte)
	}
}

func (o Objekte) Equals(o1 *Objekte) bool {
	if !o.Sha.Equals(o1.Sha) {
		return false
	}

	if !o.Akte.Equals(&o1.Akte) {
		return false
	}

	return true
}

func (o Objekte) GetGattung() gattung.Gattung {
	return gattung.Typ
}

func (o Objekte) GetAkteSha() sha.Sha {
	return o.Sha
}

func (o *Objekte) SetAkteSha(v sha.Sha) {
	o.Sha = v
}
