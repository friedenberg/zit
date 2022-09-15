package zettel_printer

import (
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/bravo/paper"
)

func (p *Printer) Sha(sh sha.Sha) (pa *paper.Paper) {
	pa = p.MakePaper()

	var sha string

	if sha, pa.Err = p.AbbreviateSha(sh); !pa.IsEmpty() {
		pa.Wrap()
		return
	}

	pa.WriteString(sha)

	return
}
