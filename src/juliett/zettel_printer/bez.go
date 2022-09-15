package zettel_printer

import (
	"github.com/friedenberg/zit/src/delta/zettel"
	"github.com/friedenberg/zit/src/paper"
)

func (p *Printer) Bezeichnung(z zettel.Zettel) (pa *paper.Paper) {
	pa = p.MakePaper()

	b := z.Bezeichnung

	if len(b) > 66 {
		b = b[:66] + "…"
	}

	pa.WriteFormat("\"%s\"", b)

	return
}
