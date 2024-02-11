package objekte_collections

import (
	"code.linenisgreat.com/zit-go/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit-go/src/charlie/collections_value"
	"code.linenisgreat.com/zit-go/src/hotel/sku"
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

	return mwk.GetKennung().String()
}

func MakeMutableSetMetadateiWithKennung() MutableSetMetadateiWithKennung {
	return collections_value.MakeMutableValueSet[*sku.Transacted](
		SkuGetKeyKeyer{},
	)
}
