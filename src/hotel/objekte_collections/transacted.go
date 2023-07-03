package objekte_collections

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/golf/objekte"
)

type MutableSetTransacted = schnittstellen.MutableSet[objekte.TransactedLike]

func MakeMutableSetTransactedUnique(c int) MutableSetTransacted {
	return collections.MakeMutableSet(
		func(sz objekte.TransactedLike) string {
			if sz == nil {
				return ""
			}
			sk := sz.GetSkuLike()

			return collections.MakeKey(
				sk.GetTai(),
				sk.GetId(),
				sk.GetObjekteSha(),
			)
		},
	)
}
