package zettel_printer

import (
	"fmt"

	"github.com/friedenberg/zit/src/bravo/zk_types"
	"github.com/friedenberg/zit/src/golf/zettel_external"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
	"github.com/friedenberg/zit/src/hotel/zettel_checked_out"
	"github.com/friedenberg/zit/src/bravo/paper"
)

func (p *Printer) ZettelCheckedOut(zco zettel_checked_out.Zettel) (pa *paper.Paper) {
	pa = p.MakePaper()

	switch zco.State {
	default:
		pa.WriteFormat("%s (unknown)", p.ZettelExternal(zco.External))

	case zettel_checked_out.StateJustCheckedOut:
		pa.WriteFormat("%s (checked out)", p.ZettelExternal(zco.External))

	case zettel_checked_out.StateJustCheckedOutButSame:
		pa.WriteFormat("%s (already checked out)", p.ZettelExternal(zco.External))

	case zettel_checked_out.StateExistsAndSame:
		pa.WriteFormat("%s (same)", p.ZettelExternal(zco.External))

	case zettel_checked_out.StateExistsAndDifferent:
		if !zco.External.ZettelFD.IsEmpty() {
			pa.WriteFormat("%s (different)", p.ZettelExternal(zco.External))
		} else if !zco.External.AkteFD.IsEmpty() {
			pa.WriteFormat("%s (Akte different)", p.ZettelExternal(zco.External))
		} else {
			pa.WriteString(fmt.Sprintf("Error! No Path or AktePath: %v", zco.External))
		}

		fallthrough

	case zettel_checked_out.StateUntracked:
		if zco.State == zettel_checked_out.StateUntracked {
			pa.WriteFormat("%s (unrecognized)", p.ZettelExternal(zco.External))
		}

		p.appendZettelCheckedOutMatches(zco.Matches, pa, zco.External)
	}

	return
}

func (p *Printer) appendZettelCheckedOutMatches(
	m zettel_checked_out.Matches,
	pa *paper.Paper,
	ex zettel_external.Zettel,
) {
	typToCollection := map[zk_types.Type]zettel_transacted.Set{
		zk_types.TypeAkte:        m.Akten,
		zk_types.TypeBezeichnung: m.Bezeichnungen,
		zk_types.TypeZettel:      m.Zettelen,
	}

	for t, c := range typToCollection {
		if c.Len() == 1 && c.Any().Named.Stored.Zettel.Equals(ex.Named.Stored.Zettel) {
		} else if c.Len() > 1 {
			c.Each(
				func(tz zettel_transacted.Zettel) (err error) {
					pa.NewLine()
					pa.WriteFormat("\t%s (%s match)", p.ZettelNamed(tz.Named), t)

					return
				},
			)
		}
	}
}
