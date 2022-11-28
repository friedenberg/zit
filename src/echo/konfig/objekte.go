package konfig

import (
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/sha"
)

type Objekte struct {
	Sha        sha.Sha
	Akte       Compiled
	tomlKonfig tomlKonfig
}

func (o *Objekte) SetAkteSha(v sha.Sha) {
	o.Sha = v
}

func (o Objekte) AkteSha() sha.Sha {
	return o.Sha
}

func (a *Objekte) Equals(b *Objekte) bool {
	panic("TODO not implemented")
	// return false
}

func (a *Objekte) Reset(b *Objekte) {
	panic("TODO not implemented")
	// return false
}

func (c Objekte) Gattung() gattung.Gattung {
	return gattung.Konfig
}
