package objekte_collections

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/golf/sku"
)

type MutableSetMetadateiWithKennung = schnittstellen.MutableSetLike[sku.SkuLike]

func MakeMutableSetMetadateiWithKennung() MutableSetMetadateiWithKennung {
	return collections.MakeMutableSet(
		func(mwk sku.SkuLike) string {
			if mwk == nil {
				return ""
			}

			return mwk.GetKey()
		},
	)
}
