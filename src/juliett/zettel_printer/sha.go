package zettel_printer

import (
	"github.com/friedenberg/zit/src/bravo/paper"
	"github.com/friedenberg/zit/src/bravo/sha"
)

func (p *Printer) Sha(sh sha.Sha) (pa *paper.Paper) {
	pa = p.MakePaper()
	pa.WriteFrom(sha.MakePaper(sh, p.Store))

	return
}
