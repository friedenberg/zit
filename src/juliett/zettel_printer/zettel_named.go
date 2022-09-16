package zettel_printer

import (
	"github.com/friedenberg/zit/src/bravo/paper"
	"github.com/friedenberg/zit/src/foxtrot/zettel_named"
)

func (p *Printer) ZettelNamed(zn zettel_named.Zettel) (pa *paper.Paper) {
	pa = p.MakePaper()

	pa.WriteString(
		p.zettelBracketed(
			p.Hinweis(zn.Hinweis),
			p.Sha(zn.Stored.Sha),
			p.Bezeichnung(zn.Stored.Zettel),
		),
	)

	return
}
