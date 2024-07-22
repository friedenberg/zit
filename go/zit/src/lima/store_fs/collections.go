package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
)

type (
	CheckedOutSet        = interfaces.SetLike[*CheckedOut]
	CheckedOutMutableSet = interfaces.MutableSetLike[*CheckedOut]
)

func MakeCheckedOutMutableSet() CheckedOutMutableSet {
	return collections_value.MakeMutableValueSet[*CheckedOut](
		nil,
	)
}
