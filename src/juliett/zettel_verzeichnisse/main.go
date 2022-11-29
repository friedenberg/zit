package zettel_verzeichnisse

import (
	"github.com/friedenberg/zit/src/india/zettel"
)

type Zettel struct {
	Transacted zettel.Transacted
	// Etiketten               tridex.Tridex
	EtikettenExpandedSorted []string
	EtikettenSorted         []string
}

func (z *Zettel) Reset(z1 *Zettel) {
	z.Transacted.Reset()
	z.EtikettenExpandedSorted = z.EtikettenExpandedSorted[:0]
	z.EtikettenSorted = z.EtikettenSorted[:0]

	if z1 != nil {
		z.Transacted = z1.Transacted

		z.EtikettenExpandedSorted = append(
			z.EtikettenExpandedSorted,
			z1.EtikettenExpandedSorted...,
		)

		z.EtikettenSorted = append(
			z.EtikettenSorted,
			z1.EtikettenSorted...,
		)
	}
}
