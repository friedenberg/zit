package zettel_external

import (
	"github.com/friedenberg/zit/src/bravo/collections"
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

		if z.Named.Stored.Sha.IsNull() {
			return ""
		}

		return z.Named.Stored.Sha.String()
	}

	return MakeMutableSet(kf, zs...)
}

func MakeMutableSetUniqueAkte(zs ...*Zettel) MutableSet {
	kf := func(z *Zettel) string {
		if z == nil {
			return ""
		}

		if z.Named.Stored.Zettel.Akte.IsNull() {
			return ""
		}

		return z.Named.Stored.Zettel.Akte.String()
	}

	return MakeMutableSet(kf, zs...)
}
