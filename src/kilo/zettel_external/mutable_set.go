package zettel_external

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections_value"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type MutableSet = schnittstellen.MutableSetLike[sku.SkuLikeExternalPtr]

type KeyerHinweis struct{}

func (k KeyerHinweis) GetKey(z sku.SkuLikeExternalPtr) string {
	if z == nil {
		return ""
	}

	return z.GetKennungLike().String()
}

func MakeMutableSetUniqueHinweis(zs ...sku.SkuLikeExternalPtr) MutableSet {
	return collections_value.MakeMutableValueSet[sku.SkuLikeExternalPtr](
		KeyerHinweis{},
		zs...)
}

type KeyerFD struct{}

func (k KeyerFD) GetKey(z sku.SkuLikeExternalPtr) string {
	if z == nil {
		return ""
	}

	return z.String()
}

func MakeMutableSetUniqueFD(zs ...sku.SkuLikeExternalPtr) MutableSet {
	return collections_value.MakeMutableValueSet[sku.SkuLikeExternalPtr](
		KeyerFD{},
		zs...)
}

type KeyerStored struct{}

func (k KeyerStored) GetKey(z sku.SkuLikeExternalPtr) string {
	if z == nil {
		return ""
	}

	if z.GetObjekteSha().IsNull() {
		return ""
	}

	return z.GetObjekteSha().String()
}

func MakeMutableSetUniqueStored(zs ...sku.SkuLikeExternalPtr) MutableSet {
	return collections_value.MakeMutableValueSet[sku.SkuLikeExternalPtr](
		KeyerStored{},
		zs...)
}

type KeyerAkte struct{}

func (k KeyerAkte) GetKey(z sku.SkuLikeExternalPtr) string {
	if z == nil {
		return ""
	}

	sh := z.GetAkteSha()

	if sh.IsNull() {
		return ""
	}

	return sh.String()
}

func MakeMutableSetUniqueAkte(zs ...sku.SkuLikeExternalPtr) MutableSet {
	return collections_value.MakeMutableValueSet[sku.SkuLikeExternalPtr](
		KeyerAkte{},
		zs...)
}
