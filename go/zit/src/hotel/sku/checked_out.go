package sku

import (
	"fmt"

	"code.linenisgreat.com/zit/src/delta/checked_out_state"
)

type CheckedOut struct {
	Internal Transacted
	External External
	State    checked_out_state.State
	IsImport bool
	Error    error
}

func (c *CheckedOut) SetError(err error) {
	if err == nil {
		return
	}

	c.State = checked_out_state.StateError
	c.Error = err
}

func (c *CheckedOut) InternalAndExternalEqualsSansTai() bool {
	return c.External.Metadatei.EqualsSansTai(
		&c.Internal.Metadatei,
	)
}

func (c *CheckedOut) DetermineState(justCheckedOut bool) {
	if c.Internal.GetObjekteSha().IsNull() {
		c.State = checked_out_state.StateUntracked
	} else if c.Internal.Metadatei.EqualsSansTai(&c.External.Metadatei) {
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

func (a *CheckedOut) String() string {
	return fmt.Sprintf("%s %s", &a.Internal, &a.External)
}
