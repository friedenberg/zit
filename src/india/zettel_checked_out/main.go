package zettel_checked_out

import (
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/alfa/logz"
	"github.com/friedenberg/zit/src/golf/stored_zettel"
)

type CheckedOut struct {
	Internal stored_zettel.Transacted
	Matches  Matches
	External stored_zettel.External
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
		c.Matches.appendToStringBuilder(sb, c.External)

	case StateAkte:
		sb.WriteString(fmt.Sprintf("[%s %s] (Hinweis not recognized)", c.External.AktePath, c.External.Stored.Zettel.Akte))
		c.Matches.appendToStringBuilder(sb, c.External)

	case StateDoesNotExist:
		sb.WriteString(fmt.Sprintf("[%s %s] (Hinweis not recognized)", c.External.Path, c.External.Stored.Sha))
		c.Matches.appendToStringBuilder(sb, c.External)
	}

	return sb.String()
}
