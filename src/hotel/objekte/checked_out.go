package objekte

import (
	"fmt"

	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type CheckedOut[
	T Akte[T],
	T1 AktePtr[T],
	T2 kennung.KennungLike[T2],
	T3 kennung.KennungLikePtr[T2],
] struct {
	Internal Transacted[T, T1, T2, T3]
	External External[T, T1, T2, T3]
	State    CheckedOutState
}

func (c *CheckedOut[T, T1, T2, T3]) DetermineState() {
	if c.Internal.Sku.ObjekteSha.IsNull() {
		c.State = CheckedOutStateUntracked
	} else if c.Internal.Sku.Metadatei.EqualsSansTai(c.External.Sku.Metadatei) {
		c.State = CheckedOutStateExistsAndSame
	} else if c.External.Sku.ObjekteSha.IsNull() {
		c.State = CheckedOutStateEmpty
	} else {
		c.State = CheckedOutStateExistsAndDifferent
	}
}

func (co CheckedOut[T, T1, T2, T3]) GetState() CheckedOutState {
	return co.State
}

func (co CheckedOut[T, T1, T2, T3]) GetInternalLike() TransactedLike {
	return co.Internal
}

func (co CheckedOut[T, T1, T2, T3]) GetExternalLike() ExternalLike {
	return co.External
}

func (a CheckedOut[T, T1, T2, T3]) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a CheckedOut[T, T1, T2, T3]) String() string {
	return fmt.Sprintf("%s %s", a.Internal.Sku, a.External.Sku)
}

func (a CheckedOut[T, T1, T2, T3]) Equals(
	b CheckedOut[T, T1, T2, T3],
) bool {
	if !a.Internal.Equals(b.Internal) {
		return false
	}

	if !a.External.Equals(b.External) {
		return false
	}

	return true
}
