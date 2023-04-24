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
	Tai        ts.Tai
	Akte       Akte
	ObjekteSha sha.Sha
	Metadatei  metadatei.Metadatei
}

func (a Objekte) GetMetadatei() metadatei.Metadatei {
	return a.Metadatei
}

func (a *Objekte) SetMetadatei(m metadatei.Metadatei) {
	a.Metadatei = m
}

func (a Objekte) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (o Objekte) GetGattung() (g schnittstellen.Gattung) {
	g = gattung.Bestandsaufnahme

	return
}

func (o Objekte) GetAkteSha() (v schnittstellen.Sha) {
	v = o.Metadatei.AkteSha

	return
}

func (o *Objekte) SetAkteSha(v schnittstellen.Sha) {
	o.Metadatei.AkteSha = sha.Make(v)
}

func (o *Objekte) SetObjekteSha(v schnittstellen.Sha) {
	o.ObjekteSha = sha.Make(v)
}

func (a Objekte) Less(b Objekte) bool {
	return a.Tai.Less(b.Tai)
}

func (a Objekte) Equals(b Objekte) bool {
	if !a.Akte.Equals(&b.Akte) {
		return false
	}

	if !a.Metadatei.Equals(b.Metadatei) {
		return false
	}

	if !a.ObjekteSha.Equals(b.ObjekteSha) {
		return false
	}

	return true
}

func (a *Objekte) Reset() {
	a.ObjekteSha = sha.Sha{}
	a.Metadatei.Reset()
	a.Tai = ts.Tai{}
	a.Akte.Reset()
}

func (a *Objekte) ResetWith(b Objekte) {
	a.Akte.ResetWith(b.Akte)
	a.ObjekteSha = b.ObjekteSha
	a.Metadatei = b.Metadatei
	a.Tai = b.Tai
}
