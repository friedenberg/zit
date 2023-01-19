package typ

import (
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type MutableSet = collections.MutableValueSet[kennung.Typ, *kennung.Typ]

func MakeMutableSet(hs ...kennung.Typ) MutableSet {
	return MutableSet(collections.MakeMutableValueSet[kennung.Typ, *kennung.Typ](hs...))
}
