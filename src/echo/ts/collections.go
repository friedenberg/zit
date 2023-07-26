package ts

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
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
	return MutableSet(collections.MakeMutableSet[Time]((Time).String, hs...))
}
