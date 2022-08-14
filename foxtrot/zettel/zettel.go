package zettel

import (
	"strings"

	"github.com/friedenberg/zit/alfa/bezeichnung"
	"github.com/friedenberg/zit/bravo/akte_ext"
	"github.com/friedenberg/zit/charlie/sha"
	"github.com/friedenberg/zit/delta/etikett"
)

type Zettel struct {
	Akte    sha.Sha
	AkteExt akte_ext.AkteExt
	//TODO-decision should this be a special etikett with different validation
	//rules?
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
