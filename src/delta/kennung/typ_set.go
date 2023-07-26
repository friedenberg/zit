package kennung

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
)

func init() {
	collections.RegisterGob[Typ]()
}

type TypSet = schnittstellen.SetLike[Typ]

func MakeTypSet(ts ...Typ) TypSet {
	return collections.MakeSet[Typ](
		(Typ).String,
		ts...,
	)
}

type TypMutableSet = schnittstellen.MutableSetLike[Typ]

func MakeTypMutableSet(ts ...Typ) TypMutableSet {
	return collections.MakeMutableSet[Typ](
		(Typ).String,
		ts...,
	)
}
