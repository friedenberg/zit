package zettel_external

import (
	"github.com/friedenberg/zit/src/charlie/collections"
)

type MutableSet struct {
	collections.MutableSetLike[*Zettel]
}

func MakeMutableSet(kf collections.KeyFunc[*Zettel], zs ...*Zettel) MutableSet {
	return MutableSet{
		MutableSetLike: collections.MakeMutableSet[*Zettel](kf, zs...),
	}
}

func MakeMutableSetUniqueHinweis(zs ...*Zettel) MutableSet {
	kf := func(z *Zettel) string {
		if z == nil {
			return ""
		}

		return z.Sku.Kennung.String()
	}

	return MakeMutableSet(kf, zs...)
}

func MakeMutableSetUniqueFD(zs ...*Zettel) MutableSet {
	kf := func(z *Zettel) string {
		if z == nil {
			return ""
		}

		return z.String()
	}

	return MakeMutableSet(kf, zs...)
}

func MakeMutableSetUniqueStored(zs ...*Zettel) MutableSet {
	kf := func(z *Zettel) string {
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

func MakeMutableSetUniqueAkte(zs ...*Zettel) MutableSet {
	kf := func(z *Zettel) string {
		if z == nil {
			return ""
		}

		if z.Objekte.Akte.IsNull() {
			return ""
		}

		return z.Objekte.Akte.String()
	}

	return MakeMutableSet(kf, zs...)
}
