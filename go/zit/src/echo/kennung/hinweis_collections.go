package kennung

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
)

func init() {
	collections_value.RegisterGobValue[Tag](nil)
}

type (
	HinweisSet        = interfaces.SetLike[Hinweis]
	HinweisMutableSet = interfaces.MutableSetLike[Hinweis]
)

func MakeHinweisMutableSet(hs ...Hinweis) HinweisMutableSet {
	return HinweisMutableSet(
		collections_value.MakeMutableValueSet[Hinweis](nil, hs...),
	)
}
