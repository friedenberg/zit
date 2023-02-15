package kennung

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
)

func init() {
	collections.RegisterGob[Hinweis]()
}

type (
	HinweisSet        = schnittstellen.Set[Hinweis]
	HinweisMutableSet = schnittstellen.MutableSet[Hinweis]
)

func MakeHinweisMutableSet(hs ...Hinweis) HinweisMutableSet {
	return HinweisMutableSet(collections.MakeMutableSet[Hinweis]((Hinweis).String, hs...))
}
