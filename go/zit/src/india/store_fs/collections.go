package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
)

type (
	CheckedOutSet        = schnittstellen.SetLike[*CheckedOut]
	CheckedOutMutableSet = schnittstellen.MutableSetLike[*CheckedOut]
)

func MakeCheckedOutMutableSet() CheckedOutMutableSet {
	return collections_value.MakeMutableValueSet[*CheckedOut](
		nil,
		// KennungKeyer[CheckedOut, *CheckedOut]{},
	)
}
