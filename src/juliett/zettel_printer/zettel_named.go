package zettel_printer

import (
	"github.com/friedenberg/zit/src/foxtrot/zettel_named"
)

func (p *Printer) ZettelNamed(zn zettel_named.Zettel) (pa *Paper) {
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
