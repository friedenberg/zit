package typ_toml

import (
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/sha"
)

type Objekte struct {
	Sha  sha.Sha
	Akte Typ
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

func (o Objekte) Gattung() gattung.Gattung {
	return gattung.Typ
}

func (o Objekte) AkteSha() sha.Sha {
	return o.Sha
}

func (o *Objekte) SetAkteSha(v sha.Sha) {
	o.Sha = v
}
