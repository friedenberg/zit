package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type (
	CheckedOutSet        = interfaces.SetLike[*sku.CheckedOut]
	CheckedOutMutableSet = interfaces.MutableSetLike[*sku.CheckedOut]
)

func MakeCheckedOutMutableSet() CheckedOutMutableSet {
	return collections_value.MakeMutableValueSet[*sku.CheckedOut](
		nil,
	)
}
