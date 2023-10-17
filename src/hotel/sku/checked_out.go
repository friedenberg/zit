package sku

import (
	"fmt"

	"github.com/friedenberg/zit/src/delta/checked_out_state"
)

type CheckedOut struct {
	Internal Transacted
	External External
	State    checked_out_state.State
}

func (c *CheckedOut) InternalAndExternalEqualsSansTai() bool {
	return c.External.Metadatei.EqualsSansTai(
		c.Internal.Metadatei,
	)
}

func (c *CheckedOut) DetermineState(justCheckedOut bool) {
	if c.Internal.GetObjekteSha().IsNull() {
		c.State = checked_out_state.StateUntracked
	} else if c.Internal.GetMetadatei().EqualsSansTai(c.External.GetMetadatei()) {
		if justCheckedOut {
			c.State = checked_out_state.StateJustCheckedOut
		} else {
			c.State = checked_out_state.StateExistsAndSame
		}
	} else {
		if justCheckedOut {
			c.State = checked_out_state.StateJustCheckedOutButDifferent
		} else {
			c.State = checked_out_state.StateExistsAndDifferent
		}
	}
}

func (a CheckedOut) String() string {
	return fmt.Sprintf("%s %s", a.Internal, a.External)
}
