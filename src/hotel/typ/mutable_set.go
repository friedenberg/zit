package typ

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections_value"
	"github.com/friedenberg/zit/src/echo/kennung"
)

type MutableSet = schnittstellen.MutableSetLike[kennung.Typ]

func MakeMutableSet(hs ...kennung.Typ) MutableSet {
	return collections_value.MakeMutableValueSet[kennung.Typ](
		nil,
		hs...,
	)
}
