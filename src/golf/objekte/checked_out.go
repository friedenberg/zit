package objekte

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

type CheckedOutState int

const (
	CheckedOutStateNotCheckedOut = CheckedOutState(iota)
	CheckedOutStateEmpty
	CheckedOutStateJustCheckedOut
	CheckedOutStateJustCheckedOutButSame
	CheckedOutStateExistsAndSame
	CheckedOutStateExistsAndDifferent
	CheckedOutStateUntracked
)

type CheckedOut[
	T schnittstellen.Objekte[T],
	T1 schnittstellen.ObjektePtr[T],
	T2 schnittstellen.Id[T2],
	T3 schnittstellen.IdPtr[T2],
	T4 any,
	T5 schnittstellen.VerzeichnissePtr[T4, T],
] struct {
	Internal Transacted[T, T1, T2, T3, T4, T5]
	External External[T, T1, T2, T3]
	State    CheckedOutState
}

func (c *CheckedOut[T, T1, T2, T3, T4, T5]) DetermineState() {
	if c.Internal.Sku.ObjekteSha.IsNull() {
		// if c.External.ExternalPathAndSha() == "" {
		// 	c.State = objekte.CheckedOutStateEmpty
		// } else {
		// 	c.State = objekte.CheckedOutStateUntracked
		// }
	} else if c.Internal.Sku.ObjekteSha.Equals(c.External.Sku.ObjekteSha) {
		c.State = CheckedOutStateExistsAndSame
	} else if c.External.Sku.ObjekteSha.IsNull() {
		c.State = CheckedOutStateEmpty
	} else {
		c.State = CheckedOutStateExistsAndDifferent
	}
}

func (co CheckedOut[T, T1, T2, T3, T4, T5]) GetState() CheckedOutState {
	return co.State
}

func (co CheckedOut[T, T1, T2, T3, T4, T5]) GetInternal() TransactedLike {
	return co.Internal
}

func (co CheckedOut[T, T1, T2, T3, T4, T5]) GetExternal() ExternalLike {
	return co.External
}
