package zettel_printer

import "github.com/friedenberg/zit/src/bravo/sha"

func (p *Printer) Sha(sh sha.Sha) (pa *Paper) {
	pa = p.MakePaper()

	var sha string

	if sha, pa.Err = p.AbbreviateSha(sh); !pa.IsEmpty() {
		pa.Wrap()
		return
	}

	pa.WriteString(sha)

	return
}
