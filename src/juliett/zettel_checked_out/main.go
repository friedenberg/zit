package zettel_checked_out

import (
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/hotel/zettel_external"
	"github.com/friedenberg/zit/src/hotel/zettel_transacted"
)

type Zettel struct {
	Internal zettel_transacted.Zettel
	Matches  Matches
	External zettel_external.Zettel
	State
}

func (c *Zettel) DetermineState() {
	if c.Internal.Named.Stored.Sha.IsNull() {
		if c.External.ExternalPathAndSha() == "" {
			c.State = StateEmpty
		} else {
			c.State = StateUntracked
		}
	} else if c.Internal.Named.Stored.Zettel.Equals(c.External.Named.Stored.Zettel) {
		c.State = StateExistsAndSame
	} else if c.External.Named.Stored.Sha.IsNull() {
		c.State = StateEmpty
	} else {
		c.State = StateExistsAndDifferent
	}
}

func (c Zettel) String() string {
	errors.PrintDebug(c)
	sb := &strings.Builder{}

	switch c.State {
	default:
		sb.WriteString(fmt.Sprintf("%v (unknown)", c.External))

	case StateExistsAndSame:
		sb.WriteString(fmt.Sprintf("%s (same)", c.External))

	case StateExistsAndDifferent:
		if !c.External.ZettelFD.IsEmpty() {
			sb.WriteString(fmt.Sprintf("%s (different)", c.External))
		} else if !c.External.AkteFD.IsEmpty() {
			sb.WriteString(fmt.Sprintf("%s (akte different)", c.External))
		} else {
			sb.WriteString(fmt.Sprintf("Error! No Path or AktePath: %v", c.External))
		}

		fallthrough

	case StateUntracked:
		if c.State == StateUntracked {
			sb.WriteString(fmt.Sprintf("%s (untracked)", c.External))
		}

		c.Matches.appendToStringBuilder(sb, c.External)
	}

	return sb.String()
}
