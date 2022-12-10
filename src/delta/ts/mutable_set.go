package ts

import "github.com/friedenberg/zit/src/delta/collections"

type MutableSet = collections.MutableValueSet[Time, *Time]

func MakeMutableSet(hs ...Time) MutableSet {
	return MutableSet(collections.MakeMutableValueSet[Time, *Time](hs...))
}
