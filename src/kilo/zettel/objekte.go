package zettel

import (
	"strings"

	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/echo/bezeichnung"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/foxtrot/kennung"
)

type Objekte struct {
	Akte          sha.Sha
	Typ           kennung.Typ
	Bezeichnung   bezeichnung.Bezeichnung
	Etiketten     kennung.EtikettSet
	Verzeichnisse Verzeichnisse2
}

func (z Objekte) Gattung() gattung.Gattung {
	return gattung.Zettel
}

func (z Objekte) AkteSha() sha.Sha {
	return z.Akte
}

func (z *Objekte) SetAkteSha(v sha.Sha) {
	z.Akte = v
}

// TODO-P4 figure out why this doesn't always work for `status`
func (z Objekte) Equals(z1 *Objekte) bool {
	if z1 == nil {
		return false
	}

	if !z.Akte.Equals(z1.Akte) {
		return false
	}

	if !z.Typ.Equals(&z1.Typ) {
		return false
	}

	if !z.Bezeichnung.Equals(z1.Bezeichnung) {
		return false
	}

	if !z.Etiketten.Equals(z1.Etiketten) {
		return false
	}

	return true
}

func (z Objekte) IsEmpty() bool {
	if strings.TrimSpace(z.Bezeichnung.String()) != "" {
		return false
	}

	if z.Etiketten.Len() > 0 {
		return false
	}

	if !z.Akte.IsNull() {
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
		z.Etiketten = z1.Etiketten.Copy()
	}
}

func (z Objekte) Description() (d string) {
	d = z.Bezeichnung.String()

	if strings.TrimSpace(d) == "" {
		d = z.Etiketten.Description()
	}

	return
}