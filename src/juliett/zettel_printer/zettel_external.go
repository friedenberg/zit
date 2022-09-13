package zettel_printer

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/golf/zettel_external"
)

func (p *Printer) ZettelExternal(ze zettel_external.Zettel) (pa *Paper) {
	pa = p.MakePaper()

	switch {
	case !ze.ZettelFD.IsEmpty():
		pa.WriteFormat("[%s %s %s]", ze.ZettelFD.Path, p.Sha(ze.Named.Stored.Sha), p.Bezeichnung(ze.Named.Stored.Zettel))

	case !ze.AkteFD.IsEmpty():
		pa.WriteFormat("[%s %s %s]", ze.AkteFD.Path, p.Sha(ze.Named.Stored.Zettel.Akte), p.Bezeichnung(ze.Named.Stored.Zettel))

	default:
		pa.Err = errors.Errorf("zettel external in unknown state: %q", ze)
	}

	return
}
