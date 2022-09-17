package zettel_printer

import (
	"github.com/friedenberg/zit/src/bravo/paper"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
)

func (p *Printer) ZettelTransacted(zt zettel_transacted.Zettel) (pa *paper.Paper) {
	pa = p.MakePaper()

	verb := ""

	switch {
	case zt.IsNew():
		verb = "created"

	default:
		verb = "updated"
	}

	pa.WriteFormat("%s (%s)", p.ZettelNamed(zt.Named), verb)

	return
}
