package objekte

import (
	"fmt"

	"github.com/friedenberg/zit/src/delta/checked_out_state"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type CheckedOut2 struct {
	Internal sku.SkuLikePtr
	External sku.SkuLikeExternalPtr
	State    checked_out_state.State
}

func (c *CheckedOut2) DetermineState(justCheckedOut bool) {
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

func (co CheckedOut2) GetState() checked_out_state.State {
	return co.State
}

func (co *CheckedOut2) GetInternalLike() sku.SkuLikePtr {
	return co.Internal
}

func (co *CheckedOut2) GetExternalLike() ExternalLike {
	return co.External
}

func (co *CheckedOut2) GetExternalLikePtr() ExternalLikePtr {
	return co.External
}

// func (a CheckedOut2) EqualsAny(b any) bool {
// 	return values.Equals(a, b)
// }

func (a CheckedOut2) String() string {
	return fmt.Sprintf("%s %s", a.Internal, a.External)
}

// func (a CheckedOut2) Equals(
// 	b CheckedOut2,
// ) bool {
// 	if !a.Internal.Equals(b.Internal) {
// 		return false
// 	}

// 	if !a.External.Equals(b.External) {
// 		return false
// 	}

// 	return true
// }
