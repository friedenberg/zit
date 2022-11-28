package typ

import (
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/fd"
	"github.com/friedenberg/zit/src/foxtrot/objekte"
)

type Objekte struct {
	Sha  sha.Sha
	Akte Akte
}

type Transacted = objekte.Transacted2[Objekte, *Objekte, kennung.Typ, *kennung.Typ]

type External struct {
	Objekte Objekte
	Kennung kennung.Typ
	Sha     sha.Sha
	FD      fd.FD
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
