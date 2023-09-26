package objekte_collections

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections_value"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type MutableSetMetadateiWithKennung = schnittstellen.MutableSetLike[*sku.Transacted]

type SkuGetKeyKeyer struct{}

func (kk SkuGetKeyKeyer) GetKey(mwk *sku.Transacted) string {
	if mwk == nil {
		return ""
	}

	return mwk.GetKey()
}

type KennungKeyer struct{}

func (kk KennungKeyer) GetKey(mwk *sku.Transacted) string {
	if mwk == nil {
		return ""
	}

	return mwk.GetKennungLike().String()
}

func MakeMutableSetMetadateiWithKennung() MutableSetMetadateiWithKennung {
	return collections_value.MakeMutableValueSet[*sku.Transacted](
		SkuGetKeyKeyer{},
	)
}
