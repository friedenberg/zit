package objekte_collections

import (
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type MutableSetMetadateiWithKennung = schnittstellen.MutableSetLike[*sku.Transacted]

// TODO remove
var MakeMutableSetMetadateiWithKennung = sku.MakeTransactedMutableSetKennung
