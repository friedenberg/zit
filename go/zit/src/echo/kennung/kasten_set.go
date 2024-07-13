package kennung

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
)

func init() {
	collections_value.RegisterGobValue[Kasten](nil)
}

type (
	KastenSet        = interfaces.SetLike[Kasten]
	KastenMutableSet = interfaces.MutableSetLike[Kasten]
)

func MakeKastenSet(ts ...Kasten) KastenSet {
	return collections_value.MakeValueSet[Kasten](
		nil,
		ts...,
	)
}

func MakeKastenMutableSet(ts ...Kasten) KastenMutableSet {
	return collections_value.MakeMutableValueSet[Kasten](
		nil,
		ts...,
	)
}
