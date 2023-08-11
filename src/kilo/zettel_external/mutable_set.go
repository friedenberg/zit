package zettel_external

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections_value"
	"github.com/friedenberg/zit/src/juliett/zettel"
)

type MutableSet = schnittstellen.MutableSetLike[*zettel.External]

type KeyerHinweis struct{}

func (k KeyerHinweis) GetKey(z *zettel.External) string {
	if z == nil {
		return ""
	}

	return z.Sku.GetKennung().String()
}

func MakeMutableSetUniqueHinweis(zs ...*zettel.External) MutableSet {
	return collections_value.MakeMutableValueSet[*zettel.External](
		KeyerHinweis{},
		zs...)
}

type KeyerFD struct{}

func (k KeyerFD) GetKey(z *zettel.External) string {
	if z == nil {
		return ""
	}

	return z.String()
}

func MakeMutableSetUniqueFD(zs ...*zettel.External) MutableSet {
	return collections_value.MakeMutableValueSet[*zettel.External](
		KeyerFD{},
		zs...)
}

type KeyerStored struct{}

func (k KeyerStored) GetKey(z *zettel.External) string {
	if z == nil {
		return ""
	}

	if z.Sku.ObjekteSha.IsNull() {
		return ""
	}

	return z.Sku.ObjekteSha.String()
}

func MakeMutableSetUniqueStored(zs ...*zettel.External) MutableSet {
	return collections_value.MakeMutableValueSet[*zettel.External](
		KeyerStored{},
		zs...)
}

type KeyerAkte struct{}

func (k KeyerAkte) GetKey(z *zettel.External) string {
	if z == nil {
		return ""
	}

	sh := z.GetAkteSha()

	if sh.IsNull() {
		return ""
	}

	return sh.String()
}

func MakeMutableSetUniqueAkte(zs ...*zettel.External) MutableSet {
	return collections_value.MakeMutableValueSet[*zettel.External](
		KeyerAkte{},
		zs...)
}
