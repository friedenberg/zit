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
}

func (co CheckedOut[T, T1, T2, T3, T4, T5]) GetInternal() TransactedLike {
	return co.Internal
}

func (co CheckedOut[T, T1, T2, T3, T4, T5]) GetExternal() ExternalLike {
	return co.External
}
