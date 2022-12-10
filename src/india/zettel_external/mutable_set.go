package zettel_external

import (
	"github.com/friedenberg/zit/src/delta/collections"
)

type MutableSet struct {
	collections.MutableSetLike[*Zettel]
}

func MakeMutableSet(kf collections.KeyFunc[*Zettel], zs ...*Zettel) MutableSet {
	return MutableSet{
		MutableSetLike: collections.MakeMutableSet[*Zettel](kf, zs...),
	}
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

		if z.Sku.Sha.IsNull() {
			return ""
		}

		return z.Sku.Sha.String()
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
