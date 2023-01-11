package bestandsaufnahme

import (
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/echo/sha"
)

type Objekte struct {
	Akte    Akte
	AkteSha sha.Sha
}

func (o Objekte) GetGattung() (g gattung.Gattung) {
	g = gattung.Bestandsaufnahme

	return
}

func (o Objekte) GetAkteSha() (v sha.Sha) {
	v = o.AkteSha

	return
}

func (o *Objekte) SetAkteSha(v sha.Sha) {
	o.AkteSha = v
}

func (a Objekte) Equals(b *Objekte) bool {
	if b == nil {
		return false
	}

	if !a.Akte.Equals(&b.Akte) {
		return false
	}

	if !a.AkteSha.Equals(b.AkteSha) {
		return false
	}

	return true
}

func (a *Objekte) Reset(b *Objekte) {
	a.Akte.Reset(&b.Akte)

	if b == nil {
		a.AkteSha = sha.Sha{}
	} else {
		a.AkteSha = b.AkteSha
	}
}
