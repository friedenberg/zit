package chrome

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (c *Chrome) QueryCheckedOut(
	qg sku.Queryable,
	f schnittstellen.FuncIter[sku.CheckedOutLike],
) (err error) {
	return
}
