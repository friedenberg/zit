package zettel_printer

import "github.com/friedenberg/zit/src/charlie/hinweis"

func (p *Printer) Hinweis(h hinweis.Hinweis) (pa *Paper) {
	pa = p.MakePaper()

	var hin string

	if p.ShouldAbbreviateHinweisen {
		if hin, pa.Err = p.AbbreviateHinweis(h); !pa.IsEmpty() {
			pa.Wrap()
			return
		}
	} else {
		hin = h.String()
	}

	pa.WriteString(hin)

	return
}
