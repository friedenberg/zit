package zettel_printer

import (
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/bravo/paper"
)

func (p *Printer) Hinweis(h hinweis.Hinweis) (pa *paper.Paper) {
	pa = p.MakePaper()

	if p.abbreviateHinweisen {
		if h, pa.Err = p.AbbreviateHinweis(h); !pa.IsEmpty() {
			pa.Wrap()
			return
		}
	}

	pa.WriteString(h.String())

	return
}
