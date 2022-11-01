package hinweis

import collections "github.com/friedenberg/zit/src/bravo/collections"

type MutableSet = collections.ValueMutableSet[Hinweis, *Hinweis]

func MakeMutableSet(hs ...Hinweis) MutableSet {
	return MutableSet(collections.MakeMutableSet[Hinweis, *Hinweis](hs...))
}
