package kennung

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
)

func init() {
	collections.RegisterGob[FD]()
}

type (
	FDSet        = schnittstellen.Set[FD]
	MutableFDSet = schnittstellen.MutableSet[FD]
)

func MakeFDSet(ts ...FD) FDSet {
	return collections.MakeSet[FD](
		(FD).String,
		ts...,
	)
}

func MakeMutableFDSet(ts ...FD) MutableFDSet {
	return collections.MakeMutableSet[FD](
		(FD).String,
		ts...,
	)
}
