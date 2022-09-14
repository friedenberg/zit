package zettel_printer

import (
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/konfig"
	"github.com/friedenberg/zit/src/hotel/store_objekten"
)

type Printer struct {
	file *os.File

	abbreviateHinweisen  bool
	abbreviateShas       bool
	newZettelShaSyntax   bool
	includeBezeichnungen bool

	printer
	*errors.Ctx
}

type printer struct {
	*store_objekten.Store
}

func Make(k konfig.Konfig, store *store_objekten.Store, f *os.File) (p *Printer) {
	p = &Printer{
		abbreviateHinweisen:  k.PrintAbbreviatedHinweisen,
		abbreviateShas:       k.PrintAbbreviatedShas,
		newZettelShaSyntax:   k.PrintNewShaSyntax,
		includeBezeichnungen: k.PrintIncludeBezeichnungen,
		printer: printer{
			Store: store,
		},
		Ctx:  &errors.Ctx{},
		file: f,
	}

	return
}
