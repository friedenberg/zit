package zettel_verzeichnisse

import (
	"github.com/friedenberg/zit/src/foxtrot/kennung"
	"github.com/friedenberg/zit/src/kilo/zettel"
)

type Zettel struct {
	Transacted zettel.Transacted
	// Etiketten               tridex.Tridex
	EtikettenExpandedSorted []string
	EtikettenSorted         []string
}

func (z *Zettel) ResetWithTransacted(z1 *zettel.Transacted) {
	if z1 != nil {
		z.Transacted.Reset(z1)
		z.EtikettenExpandedSorted = kennung.Expanded(z1.Objekte.Etiketten).SortedString()
		z.EtikettenSorted = z1.Objekte.Etiketten.SortedString()
	} else {
		z.Transacted.Reset(nil)
		z.EtikettenExpandedSorted = []string{}
		z.EtikettenSorted = []string{}
	}
}

func (z *Zettel) Reset(z1 *Zettel) {
	z.EtikettenExpandedSorted = z.EtikettenExpandedSorted[:0]
	z.EtikettenSorted = z.EtikettenSorted[:0]

	if z1 != nil {
		z.Transacted.Reset(&z1.Transacted)

		z.EtikettenExpandedSorted = append(
			z.EtikettenExpandedSorted,
			z1.EtikettenExpandedSorted...,
		)

		z.EtikettenSorted = append(
			z.EtikettenSorted,
			z1.EtikettenSorted...,
		)
	} else {
		z.Transacted.Reset(nil)
	}
}
