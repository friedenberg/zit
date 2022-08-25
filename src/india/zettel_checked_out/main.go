package zettel_checked_out

import (
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/alfa/logz"
	"github.com/friedenberg/zit/src/golf/stored_zettel"
	"github.com/friedenberg/zit/src/hotel/collections"
)

type CheckedOut struct {
	Internal      stored_zettel.Transacted
	ZettelMatches collections.SetTransacted
	AkteMatches   collections.SetTransacted
	External      stored_zettel.External
	State
}

func (c *CheckedOut) DetermineState() {
	if c.Internal.Schwanz.Time.IsZero() {
		c.State = StateDoesNotExist
	} else if c.Internal.Named.Stored.Zettel.Equals(c.External.Named.Stored.Zettel) {
		c.State = StateExistsAndSame
	} else {
		c.State = StateExistsAndDifferent
	}
}

func (c CheckedOut) String() string {
	logz.PrintDebug(c)
	sb := &strings.Builder{}

	switch c.State {
	default:
		sb.WriteString(fmt.Sprintf("%s (unknown)", c.Internal.Named))

	case StateExistsAndSame:
		sb.WriteString(fmt.Sprintf("%s (same)", c.Internal.Named))

	case StateExistsAndDifferent:
		sb.WriteString(fmt.Sprintf("%s (different)", c.External.Named))
		fallthrough

	case StateDoesNotExist:
		if c.State == StateDoesNotExist {
			sb.WriteString(fmt.Sprintf("[%s %s] (Hinweis not recognized)", c.External.Path, c.External.Stored.Sha))
		}

		if c.ZettelMatches.Len() == 1 && c.ZettelMatches.Any().Named.Stored.Zettel.Equals(c.External.Named.Stored.Zettel) {
		} else if c.ZettelMatches.Len() > 1 {
			c.ZettelMatches.Each(
				func(tz stored_zettel.Transacted) (err error) {
					sb.WriteString(fmt.Sprintf("\n\t%s (Zettel match)", tz.Named.Hinweis))
					return
				},
			)
		}
	}

	return sb.String()
}
