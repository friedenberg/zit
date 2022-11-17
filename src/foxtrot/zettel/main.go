package zettel

import (
	"strings"

	"github.com/friedenberg/zit/src/charlie/bezeichnung"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/etikett"
	"github.com/friedenberg/zit/src/delta/konfig"
	"github.com/friedenberg/zit/src/echo/typ"
)

type Zettel struct {
	Akte        sha.Sha
	Typ         typ.Typ
	Bezeichnung bezeichnung.Bezeichnung
	Etiketten   etikett.Set
}

// TODO-P2 figure out why this doesn't always work for `status`
func (z Zettel) Equals(z1 Zettel) bool {
	if !z.Akte.Equals(z1.Akte) {
		return false
	}

	if !z.Typ.Equals(z1.Typ) {
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
func (z *Zettel) Reset() {
	z.Akte = sha.Sha{}
	z.Typ = typ.Typ{}
	z.Bezeichnung = bezeichnung.Make("")
	z.Etiketten = etikett.MakeSet()
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

func (z Zettel) TypOrDefault() (t typ.Typ) {
	t = z.Typ

	if t.String() == "" {
		t = typ.Make("md")
	}

	return
}

func (z Zettel) IsInlineAkte(k konfig.Konfig) (isInline bool) {
	ts := z.Typ.String()
	isInline = k.Compiled.TypenInline.Contains(ts)

	if typKonfig, ok := k.Typen[ts]; ok {
		isInline = typKonfig.InlineAkte
	}

	return
}
