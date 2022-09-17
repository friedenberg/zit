package zettel_printer

import (
	"github.com/friedenberg/zit/src/bravo/paper"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
	"github.com/friedenberg/zit/src/india/store_working_directory"
)

func (p *Printer) FileUnrecognized(fu store_working_directory.File) (pa *paper.Paper) {
	pa = p.MakePaper()

	pa.WriteFormat("[%s %s] (not recognized)", fu.Path, p.Sha(fu.Sha))

	return
}

func (p *Printer) FileRecognized(
	fu store_working_directory.File,
	szt zettel_transacted.Set,
) (pa *paper.Paper) {
	pa = p.MakePaper()

	pa.WriteFormat("[%s %s] (Akte recognized)", fu.Path, p.Sha(fu.Sha))

	szt.Each(
		func(tz1 zettel_transacted.Zettel) (err error) {
			//TODO eliminate zettels marked as duplicates / hidden
			pa.WriteFormat("\t%s", p.ZettelNamed(tz1.Named))
			err = pa.Error()
			return
		},
	)

	return
}
