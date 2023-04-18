package bestandsaufnahme

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/echo/ts"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

type Objekte struct {
	Tai     ts.Tai
	Akte    Akte
	AkteSha sha.Sha
}

func (a Objekte) GetMetadatei() metadatei.Metadatei {
	return metadatei.Metadatei{
		AkteSha: a.AkteSha,
	}
}

func (a *Objekte) SetMetadatei(m metadatei.Metadatei) {
	a.AkteSha = m.AkteSha
}

func (a Objekte) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (o Objekte) GetGattung() (g schnittstellen.Gattung) {
	g = gattung.Bestandsaufnahme

	return
}

func (o Objekte) GetAkteSha() (v schnittstellen.Sha) {
	v = o.AkteSha

	return
}

func (o *Objekte) SetAkteSha(v schnittstellen.Sha) {
	o.AkteSha = sha.Make(v)
}

func (a Objekte) Less(b Objekte) bool {
	return a.Tai.Less(b.Tai)
}

func (a Objekte) Equals(b Objekte) bool {
	if !a.Akte.Equals(&b.Akte) {
		return false
	}

	if !a.AkteSha.Equals(b.AkteSha) {
		return false
	}

	return true
}

func (a *Objekte) Reset() {
	a.AkteSha = sha.Sha{}
	a.Tai = ts.Tai{}
	a.Akte.Reset()
}

func (a *Objekte) ResetWith(b Objekte) {
	a.Akte.ResetWith(b.Akte)
	a.AkteSha = b.AkteSha
	a.Tai = b.Tai
}
