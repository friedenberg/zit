package zettel_checked_out

import (
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/charlie/zk_types"
	"github.com/friedenberg/zit/src/golf/stored_zettel"
	"github.com/friedenberg/zit/src/hotel/collections"
)

type Matches struct {
	Akten, Bezeichnungen, Zettelen collections.SetTransacted
}

func (m Matches) appendToStringBuilder(sb *strings.Builder, ex stored_zettel.External) {
	typToCollection := map[zk_types.Type]*collections.SetTransacted{
		zk_types.TypeAkte:        &m.Akten,
		zk_types.TypeBezeichnung: &m.Bezeichnungen,
		zk_types.TypeZettel:      &m.Zettelen,
	}

	for t, c := range typToCollection {
		if c.Len() == 1 && c.Any().Named.Stored.Zettel.Equals(ex.Named.Stored.Zettel) {
		} else if c.Len() > 1 {
			c.Each(
				func(tz stored_zettel.Transacted) (err error) {
					sb.WriteString(fmt.Sprintf("\n\t%s (%s match)", tz.Named.Hinweis, t))
					return
				},
			)
		}
	}
}
