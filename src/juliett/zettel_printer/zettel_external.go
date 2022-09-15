package zettel_printer

import (
	"path/filepath"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/paper"
	"github.com/friedenberg/zit/src/golf/zettel_external"
)

func (p *Printer) ZettelExternal(ze zettel_external.Zettel) (pa *paper.Paper) {
	pa = p.MakePaper()

	bez := p.Bezeichnung(ze.Named.Stored.Zettel).String()

	var path, ref string

	switch {
	case !ze.ZettelFD.IsEmpty():
		path = ze.ZettelFD.Path
		ref = p.Sha(ze.Named.Stored.Sha).String()

	case !ze.AkteFD.IsEmpty():
		path = ze.AkteFD.Path
		ref = p.Sha(ze.Named.Stored.Zettel.Akte).String()

	default:
		pa.Err = errors.Errorf("zettel external in unknown state: %q", ze)
	}

	if pathRel, err := filepath.Rel(p.Cwd(), path); err == nil {
		path = pathRel
	}

	pa.WriteString(p.zettelBracketed(path, ref, bez))

	return
}
