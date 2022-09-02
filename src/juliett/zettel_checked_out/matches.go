package zettel_checked_out

import (
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/bravo/zk_types"
	"github.com/friedenberg/zit/src/hotel/zettel_external"
	"github.com/friedenberg/zit/src/hotel/zettel_transacted"
)

type Matches struct {
	Akten, Bezeichnungen, Zettelen zettel_transacted.Set
}

func (m Matches) appendToStringBuilder(sb *strings.Builder, ex zettel_external.Zettel) {
	typToCollection := map[zk_types.Type]*zettel_transacted.Set{
		zk_types.TypeAkte:        &m.Akten,
		zk_types.TypeBezeichnung: &m.Bezeichnungen,
		zk_types.TypeZettel:      &m.Zettelen,
	}

	for t, c := range typToCollection {
		if c.Len() == 1 && c.Any().Named.Stored.Zettel.Equals(ex.Named.Stored.Zettel) {
		} else if c.Len() > 1 {
			c.Each(
				func(tz zettel_transacted.Zettel) (err error) {
					sb.WriteString(fmt.Sprintf("\n\t%s (%s match)", tz.Named.Hinweis, t))
					return
				},
			)
		}
	}
}
