package zettel_printer

import (
	"github.com/friedenberg/zit/src/bravo/paper"
	"github.com/friedenberg/zit/src/foxtrot/zettel_named"
)

func (p *Printer) ZettelNamed(zn zettel_named.Zettel) (pa *paper.Paper) {
	pa = p.MakePaper()

	zi := p.MakeZettelish().
		Id(p.Hinweis(zn.Hinweis)).
		Rev(p.Sha(zn.Stored.Sha)).
		TypString(zn.Stored.Zettel.Typ.String()).
		Bez(p.Bezeichnung(zn.Stored.Zettel))

	pa.WriteString(zi.String())

	return
}
