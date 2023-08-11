package ts

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections_value"
)

func init() {
	collections_value.RegisterGobValue[Time](nil)
	collections_value.RegisterGobValue[Tai](nil)
}

type (
	Set        = schnittstellen.SetLike[Time]
	MutableSet = schnittstellen.MutableSetLike[Time]
)

func MakeMutableSet(hs ...Time) MutableSet {
	return MutableSet(collections_value.MakeMutableValueSet[Time](nil, hs...))
}
