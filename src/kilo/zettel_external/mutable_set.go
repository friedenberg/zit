package zettel_external

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections_value"
	"github.com/friedenberg/zit/src/hotel/external"
)

type MutableSet = schnittstellen.MutableSetLike[*external.Zettel]

type KeyerHinweis struct{}

func (k KeyerHinweis) GetKey(z *external.Zettel) string {
	if z == nil {
		return ""
	}

	return z.GetKennung().String()
}

func MakeMutableSetUniqueHinweis(zs ...*external.Zettel) MutableSet {
	return collections_value.MakeMutableValueSet[*external.Zettel](
		KeyerHinweis{},
		zs...)
}

type KeyerFD struct{}

func (k KeyerFD) GetKey(z *external.Zettel) string {
	if z == nil {
		return ""
	}

	return z.String()
}

func MakeMutableSetUniqueFD(zs ...*external.Zettel) MutableSet {
	return collections_value.MakeMutableValueSet[*external.Zettel](
		KeyerFD{},
		zs...)
}

type KeyerStored struct{}

func (k KeyerStored) GetKey(z *external.Zettel) string {
	if z == nil {
		return ""
	}

	if z.ObjekteSha.IsNull() {
		return ""
	}

	return z.ObjekteSha.String()
}

func MakeMutableSetUniqueStored(zs ...*external.Zettel) MutableSet {
	return collections_value.MakeMutableValueSet[*external.Zettel](
		KeyerStored{},
		zs...)
}

type KeyerAkte struct{}

func (k KeyerAkte) GetKey(z *external.Zettel) string {
	if z == nil {
		return ""
	}

	sh := z.GetAkteSha()

	if sh.IsNull() {
		return ""
	}

	return sh.String()
}

func MakeMutableSetUniqueAkte(zs ...*external.Zettel) MutableSet {
	return collections_value.MakeMutableValueSet[*external.Zettel](
		KeyerAkte{},
		zs...)
}
