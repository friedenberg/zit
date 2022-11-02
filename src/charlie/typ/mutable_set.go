package typ

import collections "github.com/friedenberg/zit/src/bravo/collections"

type MutableSet = collections.MutableValueSet[Typ, *Typ]

func MakeMutableSet(hs ...Typ) MutableSet {
	return MutableSet(collections.MakeMutableValueSet[Typ, *Typ](hs...))
}
