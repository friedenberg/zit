package zettel

import "github.com/friedenberg/zit/src/foxtrot/kennung"

type Verzeichnisse2 struct {
	wasPopulated bool
	// Etiketten               tridex.Tridex
	EtikettenExpandedSorted []string
	EtikettenSorted         []string
	//TODO-P3 add
	// Hidden bool
}

func (z *Verzeichnisse2) ResetWithObjekte(z1 *Objekte) {
	if z1 != nil {
		z.EtikettenExpandedSorted = kennung.Expanded(z1.Etiketten).SortedString()
		z.EtikettenSorted = z1.Etiketten.SortedString()
	} else {
		z.EtikettenExpandedSorted = []string{}
		z.EtikettenSorted = []string{}
	}
}

func (z *Verzeichnisse2) Reset(z1 *Verzeichnisse2) {
	z.ResetWithObjekte(nil)

	if z1 == nil {
		return
	}

	z.EtikettenExpandedSorted = append(
		z.EtikettenExpandedSorted,
		z1.EtikettenExpandedSorted...,
	)

	z.EtikettenSorted = append(
		z.EtikettenSorted,
		z1.EtikettenSorted...,
	)
}
