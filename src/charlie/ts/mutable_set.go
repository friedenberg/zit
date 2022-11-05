package ts

import "github.com/friedenberg/zit/src/bravo/collections"

type MutableSet = collections.MutableValueSet[Time, *Time]

func MakeMutableSet(hs ...Time) MutableSet {
	return MutableSet(collections.MakeMutableValueSet[Time, *Time](hs...))
}
