package hinweis

import "github.com/friedenberg/zit/src/charlie/collections"

type MutableSet = collections.MutableValueSet[Hinweis, *Hinweis]

func MakeMutableSet(hs ...Hinweis) MutableSet {
	return MutableSet(collections.MakeMutableValueSet[Hinweis, *Hinweis](hs...))
}