package zettel2

import (
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/bezeichnung"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type Objekte struct {
	Akte        sha.Sha
	Typ         kennung.Typ
	Bezeichnung bezeichnung.Bezeichnung
	Etiketten   kennung.EtikettSet
}

func (o Objekte) Gattung() gattung.Gattung {
	return gattung.Zettel
}

func (o Objekte) AkteSha() sha.Sha {
	return o.Akte
}

func (o *Objekte) SetAkteSha(v sha.Sha) {
	o.Akte = v
}

func (z *Objekte) Equals(z1 *Objekte) bool {
	if !z.Akte.Equals(z1.Akte) {
		return false
	}

	if !z.Typ.Equals(&z1.Typ) {
		return false
	}

	if z.Bezeichnung != z1.Bezeichnung {
		return false
	}

	if !z.Etiketten.Equals(z1.Etiketten) {
		return false
	}

	return true
}

func (z *Objekte) Reset(z1 *Objekte) {
	if z1 == nil {
		z.Akte = sha.Sha{}
		z.Typ = kennung.Typ{}
		z.Bezeichnung = bezeichnung.Make("")
		z.Etiketten = kennung.MakeEtikettSet()
	} else {
		z.Akte = z1.Akte
		z.Typ = z1.Typ
		z.Bezeichnung = z1.Bezeichnung
		z.Etiketten = z1.Etiketten
	}
}
