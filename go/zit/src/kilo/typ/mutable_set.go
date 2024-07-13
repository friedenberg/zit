package typ

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type MutableSet = interfaces.MutableSetLike[ids.Type]

func MakeMutableSet(hs ...ids.Type) MutableSet {
	return collections_value.MakeMutableValueSet[ids.Type](
		nil,
		hs...,
	)
}
