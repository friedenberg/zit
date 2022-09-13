package zettel_printer

import (
	"github.com/friedenberg/zit/src/foxtrot/zettel_named"
)

func (p *Printer) ZettelNamed(zn zettel_named.Zettel) (pa *Paper) {
	pa = p.MakePaper()

	pa.WriteFormat(
		"[%s %s %s]",
		p.Hinweis(zn.Hinweis),
		p.Sha(zn.Stored.Sha),
		p.Bezeichnung(zn.Stored.Zettel),
	)

	return
}
