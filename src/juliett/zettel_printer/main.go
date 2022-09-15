package zettel_printer

import (
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/paper"
	"github.com/friedenberg/zit/src/charlie/konfig"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/hotel/store_objekten"
)

type Printer struct {
	abbreviateHinweisen  bool
	abbreviateShas       bool
	newZettelShaSyntax   bool
	includeBezeichnungen bool

	printer
	*errors.Ctx
}

type printer struct {
	standort.Standort
	*os.File
	*store_objekten.Store
}

func Make(s standort.Standort, k konfig.Konfig, f *os.File) (p *Printer) {
	p = &Printer{
		abbreviateHinweisen:  k.PrintAbbreviatedHinweisen,
		abbreviateShas:       k.PrintAbbreviatedShas,
		newZettelShaSyntax:   k.PrintNewShaSyntax,
		includeBezeichnungen: k.PrintIncludeBezeichnungen,
		printer: printer{
			Standort: s,
			File:     f,
		},
		Ctx: &errors.Ctx{},
	}

	return
}

func (p *Printer) SetObjektenStore(s *store_objekten.Store) {
	p.Store = s
}

func (pr *Printer) MakePaper() (pa *paper.Paper) {
	return paper.Make(pr.printer.File, pr.Ctx)
}
