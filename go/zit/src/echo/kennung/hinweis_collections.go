package kennung

import (
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/charlie/collections_value"
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
