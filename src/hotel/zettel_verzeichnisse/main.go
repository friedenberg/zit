package zettel_verzeichnisse

import (
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
)

type Zettel struct {
	Transacted              zettel_transacted.Zettel
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
