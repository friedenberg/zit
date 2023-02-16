package zettel_checked_out

import (
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/values"
)

type Zettel struct {
	zettel.CheckedOut
}

func (sz Zettel) String() string {
	return collections.MakeKey(
		sz.Internal.Sku.Kopf,
		sz.Internal.Sku.Mutter[0],
		sz.Internal.Sku.Mutter[1],
		sz.Internal.Sku.Schwanz,
		sz.Internal.Sku.Kennung,
		sz.Internal.Sku.ObjekteSha,
	)
}

func (a Zettel) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a Zettel) Equals(b Zettel) bool {
	if !a.Internal.Equals(b.Internal) {
		return false
	}

	if !a.External.Equals(b.External) {
		return false
	}

	if a.State != b.State {
		return false
	}

	return true
}

func (c *Zettel) DetermineState() {
	if c.Internal.Sku.ObjekteSha.IsNull() {
		// if c.External.ExternalPathAndSha() == "" {
		// 	c.State = objekte.CheckedOutStateEmpty
		// } else {
		// 	c.State = objekte.CheckedOutStateUntracked
		// }
	} else if c.Internal.Sku.ObjekteSha.Equals(c.External.Sku.ObjekteSha) {
		c.State = objekte.CheckedOutStateExistsAndSame
	} else if c.External.Sku.ObjekteSha.IsNull() {
		c.State = objekte.CheckedOutStateEmpty
	} else {
		c.State = objekte.CheckedOutStateExistsAndDifferent
	}
}
