package zettel_printer

import (
	"github.com/friedenberg/zit/src/bravo/paper"
	"github.com/friedenberg/zit/src/delta/zettel"
)

func (p *Printer) Bezeichnung(z zettel.Zettel) (pa *paper.Paper) {
	pa = p.MakePaper()

	b := z.Bezeichnung.String()

	switch {
	case len(b) > 66:
		b = b[:66] + "…"

	case len(b) == 0:
		b = z.Etiketten.Description()
	}

	pa.WriteFormat("\"%s\"", b)

	return
}
