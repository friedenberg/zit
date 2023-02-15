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
	Set        = schnittstellen.Set[Time]
	MutableSet = schnittstellen.MutableSet[Time]
)

func MakeMutableSet(hs ...Time) MutableSet {
	return MutableSet(collections.MakeMutableSet[Time]((Time).String, hs...))
}
