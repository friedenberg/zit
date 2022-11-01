package typ

import collections "github.com/friedenberg/zit/src/bravo/collections"

type MutableSet = collections.ValueMutableSet[Typ, *Typ]

func MakeMutableSet(hs ...Typ) MutableSet {
	return MutableSet(collections.MakeMutableSet[Typ, *Typ](hs...))
}
