package kennung

import "github.com/friedenberg/zit/src/charlie/collections"

type FDSet = collections.ValueSet[FD, *FD]

func MakeFDSet(ts ...FD) FDSet {
	return collections.MakeValueSet[FD, *FD](
		ts...,
	)
}

type FDMutableSet = collections.MutableValueSet[FD, *FD]

func MakeFDMutableSet(ts ...FD) FDMutableSet {
	return collections.MakeMutableValueSet[FD, *FD](
		ts...,
	)
}
