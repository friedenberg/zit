package query2

import (
	"code.linenisgreat.com/zit/src/hotel/matcher_proto"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/matcher"
)

func MakeGroupFromCheckedOutSet(
	cos sku.CheckedOutSet,
) (q matcher_proto.QueryGroup, err error) {
	return matcher.MakeGroupFromCheckedOutSet(cos)
}
