package zettel_external

import collections "github.com/friedenberg/zit/src/bravo/collections"

type MutableSet struct {
	collections.MutableSetGeneric[*Zettel]
}

func MakeMutableSet(zs ...*Zettel) MutableSet {
	kf := func(z *Zettel) string {
		return z.String()
	}

	return MutableSet{
		MutableSetGeneric: collections.MakeMutableSetGeneric[*Zettel](kf, zs...),
	}
}
