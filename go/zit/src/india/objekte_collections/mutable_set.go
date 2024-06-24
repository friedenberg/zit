package objekte_collections

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type MutableSet = schnittstellen.MutableSetLike[*sku.ExternalFS]

type KeyerHinweis struct{}

func (k KeyerHinweis) GetKey(z *sku.ExternalFS) string {
	if z == nil {
		return ""
	}

	return z.GetKennung().String()
}

func MakeMutableSetUniqueHinweis(zs ...*sku.ExternalFS) MutableSet {
	return collections_value.MakeMutableValueSet[*sku.ExternalFS](
		KeyerHinweis{},
		zs...)
}

type KeyerFD struct{}

func (k KeyerFD) GetKey(z *sku.ExternalFS) string {
	if z == nil {
		return ""
	}

	return z.String()
}

func MakeMutableSetUniqueFD(zs ...*sku.ExternalFS) MutableSet {
	return collections_value.MakeMutableValueSet[*sku.ExternalFS](
		KeyerFD{},
		zs...)
}

type KeyerStored struct{}

func (k KeyerStored) GetKey(z *sku.ExternalFS) string {
	if z == nil {
		return ""
	}

	if z.GetObjekteSha().IsNull() {
		return ""
	}

	return z.GetObjekteSha().String()
}

func MakeMutableSetUniqueStored(zs ...*sku.ExternalFS) MutableSet {
	return collections_value.MakeMutableValueSet[*sku.ExternalFS](
		KeyerStored{},
		zs...)
}

type KeyerAkte struct{}

func (k KeyerAkte) GetKey(z *sku.ExternalFS) string {
	if z == nil {
		return ""
	}

	sh := z.GetAkteSha()

	if sh.IsNull() {
		return ""
	}

	return sh.String()
}

func MakeMutableSetUniqueAkte(zs ...*sku.ExternalFS) MutableSet {
	return collections_value.MakeMutableValueSet[*sku.ExternalFS](
		KeyerAkte{},
		zs...)
}
