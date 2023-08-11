package kennung

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/collections_value"
)

func init() {
	collections.RegisterGob[Kasten]()
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
	return collections.MakeMutableSet[Kasten](
		(Kasten).String,
		ts...,
	)
}
