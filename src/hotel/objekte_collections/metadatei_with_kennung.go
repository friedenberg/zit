package objekte_collections

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/foxtrot/sku"
)

type MutableSetMetadateiWithKennung = schnittstellen.MutableSet[sku.WithKennungInterface]

func MakeMutableSetMetadateiWithKennung() MutableSetMetadateiWithKennung {
	return collections.MakeMutableSet(
		func(mwk sku.WithKennungInterface) string {
			return collections.MakeKey(mwk.Kennung)
		},
	)
}
