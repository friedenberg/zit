package objekte_collections

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

type MutableSetMetadateiWithKennung = schnittstellen.MutableSet[metadatei.WithKennungInterface]

func MakeMutableSetMetadateiWithKennung() MutableSetMetadateiWithKennung {
	return collections.MakeMutableSet(
		func(mwk metadatei.WithKennungInterface) string {
			return collections.MakeKey(mwk.Kennung)
		},
	)
}
