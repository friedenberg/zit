package zettel_checked_out

import (
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/golf/stored_zettel"
	"github.com/friedenberg/zit/src/hotel/collections"
)

type CheckedOut struct {
	Internal stored_zettel.Transacted
	Matches  collections.SetTransacted
	External stored_zettel.External
}

func (c CheckedOut) String() string {
	sb := &strings.Builder{}

	if c.Internal.Zettel.Equals(c.External.Zettel) {
		sb.WriteString(fmt.Sprintf("%s (same)", c.Internal.Named))
	} else {
		sb.WriteString(fmt.Sprintf("%s (different)", c.External.Named))

		if c.Matches.Len() == 1 && c.Matches.Any().Zettel.Equals(c.External.Zettel) {
		} else if c.Matches.Len() > 1 {
			c.Matches.Each(
				func(tz stored_zettel.Transacted) (err error) {
					sb.WriteString(fmt.Sprintf("\n\t%s (match)", c.External.Hinweis))
					return
				},
			)
		}
	}

	return sb.String()
}
