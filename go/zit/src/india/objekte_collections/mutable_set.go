package objekte_collections

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type MutableSet = schnittstellen.MutableSetLike[*sku.External]

type KeyerHinweis struct{}

func (k KeyerHinweis) GetKey(z *sku.External) string {
	if z == nil {
		return ""
	}

	return z.GetKennung().String()
}

func MakeMutableSetUniqueHinweis(zs ...*sku.External) MutableSet {
	return collections_value.MakeMutableValueSet[*sku.External](
		KeyerHinweis{},
		zs...)
}

type KeyerFD struct{}

func (k KeyerFD) GetKey(z *sku.External) string {
	if z == nil {
		return ""
	}

	return z.String()
}

func MakeMutableSetUniqueFD(zs ...*sku.External) MutableSet {
	return collections_value.MakeMutableValueSet[*sku.External](
		KeyerFD{},
		zs...)
}

type KeyerStored struct{}

func (k KeyerStored) GetKey(z *sku.External) string {
	if z == nil {
		return ""
	}

	if z.GetObjekteSha().IsNull() {
		return ""
	}

	return z.GetObjekteSha().String()
}

func MakeMutableSetUniqueStored(zs ...*sku.External) MutableSet {
	return collections_value.MakeMutableValueSet[*sku.External](
		KeyerStored{},
		zs...)
}

type KeyerAkte struct{}

func (k KeyerAkte) GetKey(z *sku.External) string {
	if z == nil {
		return ""
	}

	sh := z.GetAkteSha()

	if sh.IsNull() {
		return ""
	}

	return sh.String()
}

func MakeMutableSetUniqueAkte(zs ...*sku.External) MutableSet {
	return collections_value.MakeMutableValueSet[*sku.External](
		KeyerAkte{},
		zs...)
}
