package objekte_collections

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/hotel/objekte"
)

type (
	TransactedLikeSet        schnittstellen.SetLike[objekte.TransactedLike]
	TransactedLikeMutableSet schnittstellen.MutableSetLike[objekte.TransactedLike]
)

func MakeTransactedLikeMutableSetKennung() TransactedLikeMutableSet {
	return collections.MakeMutableSet(
		func(tl objekte.TransactedLike) string {
			if tl == nil {
				return ""
			}

			return tl.GetKennungLike().String()
		},
	)
}
