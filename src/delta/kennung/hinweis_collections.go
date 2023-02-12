package kennung

import "github.com/friedenberg/zit/src/charlie/collections"

type (
	HinweisSet        = collections.ValueSet[Hinweis, *Hinweis]
	HinweisMutableSet = collections.MutableValueSet[Hinweis, *Hinweis]
)

func MakeHinweisMutableSet(hs ...Hinweis) HinweisMutableSet {
	return HinweisMutableSet(collections.MakeMutableValueSet[Hinweis, *Hinweis](hs...))
}
