package typ

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type MutableSet = schnittstellen.MutableSet[kennung.Typ]

func MakeMutableSet(hs ...kennung.Typ) MutableSet {
	return collections.MakeMutableSet[kennung.Typ](
		(kennung.Typ).String,
		hs...,
	)
}
