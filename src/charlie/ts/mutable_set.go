package ts

import collections "github.com/friedenberg/zit/src/bravo/collections"

type MutableSet = collections.ValueMutableSet[Time, *Time]

func MakeMutableSet(hs ...Time) MutableSet {
	return MutableSet(collections.MakeMutableSet[Time, *Time](hs...))
}
