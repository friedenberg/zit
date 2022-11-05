package zettel_printer

import (
	"github.com/friedenberg/zit/src/bravo/paper"
	"github.com/friedenberg/zit/src/charlie/hinweis"
)

func (p *Printer) Hinweis(h hinweis.Hinweis) (pa *paper.Paper) {
	pa = p.MakePaper()

	if p.abbreviateHinweisen {
		h, _ = p.AbbreviateHinweis(h)
	}

	pa.WriteString(h.String())

	return
}
