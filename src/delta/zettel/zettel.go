package zettel

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/bezeichnung"
	"github.com/friedenberg/zit/src/alfa/typ"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/etikett"
)

type Zettel struct {
	Akte        sha.Sha
	Typ         typ.Typ
	Bezeichnung bezeichnung.Bezeichnung
	Etiketten   etikett.Set
}

func (z Zettel) Description() (d string) {
	d = z.Bezeichnung.String()

	if strings.TrimSpace(d) == "" {
		sb := &strings.Builder{}
		first := true

		for _, e1 := range z.Etiketten.Sorted() {
			if !first {
				sb.WriteString(", ")
			}

			sb.WriteString(e1.String())

			first = false
		}

		d = sb.String()
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
		t = typ.Typ{Value: "md"}
	}

	return
}

func (z Zettel) AkteExt() (ext string) {
	t := z.TypOrDefault()
	ext = t.String()

	return
}
