package zettel

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/bezeichnung"
	"github.com/friedenberg/zit/src/schnittstellen"
)

type Objekte struct {
	Akte          sha.Sha
	Typ           kennung.Typ
	Bezeichnung   bezeichnung.Bezeichnung
	Etiketten     kennung.EtikettSet
	Verzeichnisse Verzeichnisse
}

func (z Objekte) GetGattung() schnittstellen.Gattung {
	return gattung.Zettel
}

func (z Objekte) GetAkteSha() schnittstellen.Sha {
	return z.Akte
}

func (z *Objekte) SetAkteSha(v schnittstellen.Sha) {
	z.Akte = sha.Make(v)
}

func (z Objekte) Equals(z1 Objekte) bool {
	errors.TodoP4("figure out why this doesn't always work for `status`")
	if !z.Akte.Equals(z1.Akte) {
		return false
	}

	if !z.Typ.Equals(z1.Typ) {
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
