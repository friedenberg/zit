package zettel_external

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/juliett/zettel"
)

type MutableSet struct {
	schnittstellen.MutableSet[*zettel.External]
}

func MakeMutableSet(
	kf collections.KeyFunc[*zettel.External],
	zs ...*zettel.External,
) MutableSet {
	return MutableSet{
		MutableSet: collections.MakeMutableSet[*zettel.External](kf, zs...),
	}
}

func MakeMutableSetUniqueHinweis(zs ...*zettel.External) MutableSet {
	kf := func(z *zettel.External) string {
		if z == nil {
			return ""
		}

		return z.Sku.Kennung.String()
	}

	return MakeMutableSet(kf, zs...)
}

func MakeMutableSetUniqueFD(zs ...*zettel.External) MutableSet {
	kf := func(z *zettel.External) string {
		if z == nil {
			return ""
		}

		return z.String()
	}

	return MakeMutableSet(kf, zs...)
}

func MakeMutableSetUniqueStored(zs ...*zettel.External) MutableSet {
	kf := func(z *zettel.External) string {
		if z == nil {
			return ""
		}

		if z.Sku.ObjekteSha.IsNull() {
			return ""
		}

		return z.Sku.ObjekteSha.String()
	}

	return MakeMutableSet(kf, zs...)
}

func MakeMutableSetUniqueAkte(zs ...*zettel.External) MutableSet {
	kf := func(z *zettel.External) string {
		if z == nil {
			return ""
		}

		if z.GetAkteSha().IsNull() {
			return ""
		}

		return z.GetAkteSha().String()
	}

	return MakeMutableSet(kf, zs...)
}
