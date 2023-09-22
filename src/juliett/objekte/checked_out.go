package objekte

import (
	"fmt"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/delta/checked_out_state"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type (
	CheckedOutLike interface {
		schnittstellen.ValueLike
		GetInternalLike() sku.SkuLikePtr
		GetExternalLike() ExternalLike
		GetState() checked_out_state.State
	}

	CheckedOutLikePtr interface {
		CheckedOutLike
		GetInternalLikePtr() sku.SkuLikePtr
		GetExternalLikePtr() ExternalLikePtr
		SetExternalLikePtr(ExternalLikePtr) error
		DetermineState(justCheckedOut bool)
		SetState(checked_out_state.State)
	}

	CheckedOut[
		T2 kennung.KennungLike[T2],
		T3 kennung.KennungLikePtr[T2],
	] struct {
		Internal sku.Transacted[T2, T3]
		External sku.External2
		State    checked_out_state.State
	}
)

func (c *CheckedOut[T2, T3]) DetermineState(justCheckedOut bool) {
	if c.Internal.ObjekteSha.IsNull() {
		c.State = checked_out_state.StateUntracked
	} else if c.Internal.Metadatei.EqualsSansTai(c.External.Metadatei) {
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

func (co CheckedOut[T2, T3]) GetState() checked_out_state.State {
	return co.State
}

func (co *CheckedOut[T2, T3]) SetState(v checked_out_state.State) {
	co.State = v
}

func (co *CheckedOut[T2, T3]) GetInternalLike() sku.SkuLikePtr {
	return &co.Internal
}

func (co *CheckedOut[T2, T3]) GetExternalLike() ExternalLike {
	return &co.External
}

func (co *CheckedOut[T2, T3]) GetExternalLikePtr() ExternalLikePtr {
	return &co.External
}

func (co *CheckedOut[T2, T3]) GetInternalLikePtr() sku.SkuLikePtr {
	return &co.Internal
}

func (co *CheckedOut[T2, T3]) SetExternalLikePtr(
	v ExternalLikePtr,
) (err error) {
	if err = co.External.SetFromSkuLike(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a CheckedOut[T2, T3]) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a CheckedOut[T2, T3]) String() string {
	return fmt.Sprintf("%s %s", a.Internal, a.External)
}

func (a CheckedOut[T2, T3]) Equals(
	b CheckedOut[T2, T3],
) bool {
	if !a.Internal.Equals(b.Internal) {
		return false
	}

	if !a.External.Equals(b.External) {
		return false
	}

	return true
}
