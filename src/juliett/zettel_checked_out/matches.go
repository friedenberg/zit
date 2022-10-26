package zettel_checked_out

import (
	"fmt"
	"strings"

	gattung "github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/golf/zettel_external"
	"github.com/friedenberg/zit/src/hotel/zettel_transacted"
)

type Matches struct {
	Akten, Bezeichnungen, Zettelen zettel_transacted.Set
}

func (m Matches) appendToStringBuilder(sb *strings.Builder, ex zettel_external.Zettel) {
	typToCollection := map[gattung.Gattung]*zettel_transacted.Set{
		gattung.Akte:        &m.Akten,
		gattung.Bezeichnung: &m.Bezeichnungen,
		gattung.Zettel:      &m.Zettelen,
	}

	for t, c := range typToCollection {
		if c.Len() == 1 && c.Any().Named.Stored.Zettel.Equals(ex.Named.Stored.Zettel) {
		} else if c.Len() > 1 {
			c.Each(
				func(tz *zettel_transacted.Zettel) (err error) {
					sb.WriteString(fmt.Sprintf("\n\t%s (%s match)", tz.Named, t))
					return
				},
			)
		}
	}
}
