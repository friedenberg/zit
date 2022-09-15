package zettel_printer

import (
	"github.com/friedenberg/zit/src/foxtrot/zettel_named"
	"github.com/friedenberg/zit/src/paper"
)

func (p *Printer) ZettelNamed(zn zettel_named.Zettel) (pa *paper.Paper) {
	pa = p.MakePaper()

	pa.WriteString(
		p.zettelBracketed(
			p.Hinweis(zn.Hinweis).String(),
			p.Sha(zn.Stored.Sha).String(),
			p.Bezeichnung(zn.Stored.Zettel).String(),
		),
	)

	return
}
