package zettel_verzeichnisse

import (
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
)

type Zettel struct {
	Transacted              zettel_transacted.Zettel
	EtikettenExpandedSorted []string
	EtikettenSorted         []string
}

func (z *Zettel) Reset() {
	z.Transacted.Reset()
	z.EtikettenExpandedSorted = z.EtikettenExpandedSorted[:0]
	z.EtikettenSorted = z.EtikettenSorted[:0]
}
