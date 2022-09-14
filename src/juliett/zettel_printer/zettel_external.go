package zettel_printer

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/golf/zettel_external"
)

func (p *Printer) ZettelExternal(ze zettel_external.Zettel) (pa *Paper) {
	pa = p.MakePaper()

	switch {
	case !ze.ZettelFD.IsEmpty():
		pa.WriteString(
			p.zettelBracketed(
				ze.ZettelFD.Path,
				p.Sha(ze.Named.Stored.Sha).String(),
				p.Bezeichnung(ze.Named.Stored.Zettel).String(),
			),
		)

	case !ze.AkteFD.IsEmpty():
		pa.WriteString(
			p.zettelBracketed(
				ze.AkteFD.Path,
				p.Sha(ze.Named.Stored.Zettel.Akte).String(),
				p.Bezeichnung(ze.Named.Stored.Zettel).String(),
			),
		)

	default:
		pa.Err = errors.Errorf("zettel external in unknown state: %q", ze)
	}

	return
}
