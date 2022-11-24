package zettel

import (
	"strings"

	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/bezeichnung"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/objekte"
)

type Zettel struct {
	Akte        sha.Sha
	Typ         kennung.Typ
	Bezeichnung bezeichnung.Bezeichnung
	Etiketten   kennung.EtikettSet
}

type Stored = objekte.Stored[Zettel, *Zettel]

type Named = objekte.Named[
	Zettel,
	*Zettel,
	hinweis.Hinweis,
	*hinweis.Hinweis,
]

func (z Zettel) Gattung() gattung.Gattung {
	return gattung.Zettel
}

// TODO-P2 figure out why this doesn't always work for `status`
func (z *Zettel) Equals(z1 *Zettel) bool {
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

func (z Zettel) IsEmpty() bool {
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

// TODO-P3 use reset with pointer pattern
func (z *Zettel) Reset(z1 *Zettel) {
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

func (z Zettel) Description() (d string) {
	d = z.Bezeichnung.String()

	if strings.TrimSpace(d) == "" {
		d = z.Etiketten.Description()
	}

	return
}

func (z Zettel) DescriptionAndTags() (d string) {
	sb := &strings.Builder{}

	sb.WriteString(z.Bezeichnung.String())

	for _, e1 := range z.Etiketten.Sorted() {
		sb.WriteString(", ")
		sb.WriteString(e1.String())
	}

	d = sb.String()

	return
}
