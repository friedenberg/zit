package zettel_printer

import (
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/hotel/store_objekten"
)

type Printer struct {
	*os.File
	ShouldAbbreviateHinweisen bool
	printer
	*errors.Ctx
}

type printer struct {
	*store_objekten.Store
}

func Make(store *store_objekten.Store, f *os.File) (p *Printer) {
	p = &Printer{
		printer: printer{
			Store: store,
		},
		Ctx:  &errors.Ctx{},
		File: f,
	}

	return
}
