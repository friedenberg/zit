package objekte_collections

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/golf/objekte"
)

type (
	TransactedLikeSet        schnittstellen.Set[objekte.TransactedLike]
	TransactedLikeMutableSet schnittstellen.MutableSet[objekte.TransactedLike]
)

func MakeTransactedLikeMutableSetKennung() TransactedLikeMutableSet {
	return collections.MakeMutableSet(
		func(tl objekte.TransactedLike) string {
			if tl == nil {
				return ""
			}

			return tl.GetKennung().String()
		},
	)
}
