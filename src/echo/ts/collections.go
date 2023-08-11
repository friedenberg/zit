package ts

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/collections_value"
)

func init() {
	collections.RegisterGob[Time]()
	collections.RegisterGob[Tai]()
}

type (
	Set        = schnittstellen.SetLike[Time]
	MutableSet = schnittstellen.MutableSetLike[Time]
)

func MakeMutableSet(hs ...Time) MutableSet {
	return MutableSet(collections_value.MakeMutableValueSet[Time](nil, hs...))
}
