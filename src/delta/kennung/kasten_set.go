package kennung

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
)

func init() {
	collections.RegisterGob[Kasten]()
}

type (
	KastenSet        = schnittstellen.Set[Kasten]
	KastenMutableSet = schnittstellen.MutableSet[Kasten]
)

func MakeKastenSet(ts ...Kasten) KastenSet {
	return collections.MakeSet[Kasten](
		(Kasten).String,
		ts...,
	)
}

func MakeKastenMutableSet(ts ...Kasten) KastenMutableSet {
	return collections.MakeMutableSet[Kasten](
		(Kasten).String,
		ts...,
	)
}
