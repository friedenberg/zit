package konfig

import (
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/foxtrot/objekte"
)

type Transacted = objekte.Transacted[Objekte, *Objekte, kennung.Konfig, *kennung.Konfig]

type Objekte struct {
	Sha  sha.Sha
	Akte Toml
}

func (o *Objekte) SetAkteSha(v sha.Sha) {
	o.Sha = v
}

func (o Objekte) AkteSha() sha.Sha {
	return o.Sha
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
		a.Sha = b.Sha
		a.Akte = b.Akte
	} else {
		a.Sha = sha.Sha{}
		//TODO
		// a.Akte = MakeDefaultCompiled()
	}
}

func (c Objekte) Gattung() gattung.Gattung {
	return gattung.Konfig
}
