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
	if c.Internal.Sku.ObjekteSha.IsNull() {
		if c.External.ExternalPathAndSha() == "" {
			c.State = StateEmpty
		} else {
			c.State = StateUntracked
		}
	} else if c.Internal.Sku.ObjekteSha.Equals(c.External.Sku.ObjekteSha) {
		c.State = StateExistsAndSame
	} else if c.External.Sku.ObjekteSha.IsNull() {
		c.State = StateEmpty
	} else {
		c.State = StateExistsAndDifferent
	}
}
