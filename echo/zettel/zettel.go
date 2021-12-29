package zettel

import "strings"

type Zettel struct {
	Akte    _Sha
	AkteExt _AkteExt
	//TODO make this a special etikett with different validation rules
	Bezeichnung _Bezeichnung
	Etiketten   _EtikettSet
}

func (z Zettel) Description() (d string) {
	d = z.Bezeichnung.String()

	if d == "" {
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
