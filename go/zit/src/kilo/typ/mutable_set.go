package typ

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
)

type MutableSet = interfaces.MutableSetLike[kennung.Type]

func MakeMutableSet(hs ...kennung.Type) MutableSet {
	return collections_value.MakeMutableValueSet[kennung.Type](
		nil,
		hs...,
	)
}
