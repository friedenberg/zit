package kennung

import (
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/charlie/collections_value"
)

func init() {
	collections_value.RegisterGobValue[Kasten](nil)
}

type (
	KastenSet        = schnittstellen.SetLike[Kasten]
	KastenMutableSet = schnittstellen.MutableSetLike[Kasten]
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
