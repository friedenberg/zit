package kennung

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections_value"
)

func init() {
	collections_value.RegisterGobValue[Etikett](nil)
}

type (
	HinweisSet        = schnittstellen.SetLike[Hinweis]
	HinweisMutableSet = schnittstellen.MutableSetLike[Hinweis]
)

func MakeHinweisMutableSet(hs ...Hinweis) HinweisMutableSet {
	return HinweisMutableSet(
		collections_value.MakeMutableValueSet[Hinweis](nil, hs...),
	)
}
