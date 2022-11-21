package typ

import "github.com/friedenberg/zit/src/bravo/collections"

type MutableSet = collections.MutableValueSet[Kennung, *Kennung]

func MakeMutableSet(hs ...Kennung) MutableSet {
	return MutableSet(collections.MakeMutableValueSet[Kennung, *Kennung](hs...))
}
