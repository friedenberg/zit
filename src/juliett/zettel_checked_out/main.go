package zettel_checked_out

import (
	"github.com/friedenberg/zit/src/india/zettel_external"
	"github.com/friedenberg/zit/src/kilo/zettel"
)

type Zettel struct {
	Internal zettel.Transacted
	External zettel_external.Zettel
	State
}

func (c *Zettel) DetermineState() {
	if c.Internal.Sku.Sha.IsNull() {
		if c.External.ExternalPathAndSha() == "" {
			c.State = StateEmpty
		} else {
			c.State = StateUntracked
		}
	} else if c.Internal.Sku.Sha.Equals(c.External.Sku.Sha) {
		c.State = StateExistsAndSame
	} else if c.External.Sku.Sha.IsNull() {
		c.State = StateEmpty
	} else {
		c.State = StateExistsAndDifferent
	}
}
