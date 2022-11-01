package zettel_named

import "github.com/friedenberg/zit/src/bravo/collections"

type MutableSet struct {
	collections.MutableSetLike[*Zettel]
}

func MakeMutableSet(kf collections.KeyFunc[*Zettel], zs ...*Zettel) MutableSet {
	return MutableSet{
		MutableSetLike: collections.MakeMutableSetGeneric[*Zettel](kf, zs...),
	}
}

func MakeMutableSetUniqueHinweisen(zs ...*Zettel) MutableSet {
	kf := func(z *Zettel) string {
		return z.Hinweis.String()
	}

	return MakeMutableSet(kf, zs...)
}
