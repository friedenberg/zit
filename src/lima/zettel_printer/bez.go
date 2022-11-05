package zettel_printer

import (
	"github.com/friedenberg/zit/src/alfa/bezeichnung"
	"github.com/friedenberg/zit/src/bravo/paper"
	"github.com/friedenberg/zit/src/delta/zettel"
)

func (p *Printer) Bezeichnung(z zettel.Zettel) (pa *paper.Paper) {
	pa = p.MakePaper()

	b := &z.Bezeichnung

	if b.IsEmpty() {
		pa.WriteFormat("\"%s\"", z.Etiketten.Description())
	} else {
		f := bezeichnung.MakeCliFormat()
		f(pa, b)
	}

	return
}
