package zettel_checked_out

import (
	"github.com/friedenberg/zit/src/india/zettel"
	"github.com/friedenberg/zit/src/india/zettel_external"
)

type Zettel struct {
	Internal zettel.Transacted
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
