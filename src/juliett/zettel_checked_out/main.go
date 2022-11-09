package zettel_checked_out

import (
	"github.com/friedenberg/zit/src/golf/zettel_external"
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
	} else if c.Internal.Named.Stored.Sha.Equals(c.External.Named.Stored.Sha) {
		c.State = StateExistsAndSame
	} else if c.External.Named.Stored.Sha.IsNull() {
		c.State = StateEmpty
	} else {
		c.State = StateExistsAndDifferent
	}
}
