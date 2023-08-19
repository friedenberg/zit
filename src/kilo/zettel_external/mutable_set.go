package zettel_external

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections_value"
	"github.com/friedenberg/zit/src/golf/sku"
)

type MutableSet = schnittstellen.MutableSetLike[*sku.ExternalZettel]

type KeyerHinweis struct{}

func (k KeyerHinweis) GetKey(z *sku.ExternalZettel) string {
	if z == nil {
		return ""
	}

	return z.GetKennung().String()
}

func MakeMutableSetUniqueHinweis(zs ...*sku.ExternalZettel) MutableSet {
	return collections_value.MakeMutableValueSet[*sku.ExternalZettel](
		KeyerHinweis{},
		zs...)
}

type KeyerFD struct{}

func (k KeyerFD) GetKey(z *sku.ExternalZettel) string {
	if z == nil {
		return ""
	}

	return z.String()
}

func MakeMutableSetUniqueFD(zs ...*sku.ExternalZettel) MutableSet {
	return collections_value.MakeMutableValueSet[*sku.ExternalZettel](
		KeyerFD{},
		zs...)
}

type KeyerStored struct{}

func (k KeyerStored) GetKey(z *sku.ExternalZettel) string {
	if z == nil {
		return ""
	}

	if z.ObjekteSha.IsNull() {
		return ""
	}

	return z.ObjekteSha.String()
}

func MakeMutableSetUniqueStored(zs ...*sku.ExternalZettel) MutableSet {
	return collections_value.MakeMutableValueSet[*sku.ExternalZettel](
		KeyerStored{},
		zs...)
}

type KeyerAkte struct{}

func (k KeyerAkte) GetKey(z *sku.ExternalZettel) string {
	if z == nil {
		return ""
	}

	sh := z.GetAkteSha()

	if sh.IsNull() {
		return ""
	}

	return sh.String()
}

func MakeMutableSetUniqueAkte(zs ...*sku.ExternalZettel) MutableSet {
	return collections_value.MakeMutableValueSet[*sku.ExternalZettel](
		KeyerAkte{},
		zs...)
}
