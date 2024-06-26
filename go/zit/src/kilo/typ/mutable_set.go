package typ

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
)

type MutableSet = schnittstellen.MutableSetLike[kennung.Typ]

func MakeMutableSet(hs ...kennung.Typ) MutableSet {
	return collections_value.MakeMutableValueSet[kennung.Typ](
		nil,
		hs...,
	)
}
