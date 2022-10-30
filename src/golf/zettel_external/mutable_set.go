package zettel_external

import (
	collections "github.com/friedenberg/zit/src/bravo/collections"
)

type MutableSet struct {
	collections.MutableSetGeneric[*Zettel]
}

func MakeMutableSet(kf collections.KeyFunc[*Zettel], zs ...*Zettel) MutableSet {
	return MutableSet{
		MutableSetGeneric: collections.MakeMutableSetGeneric[*Zettel](kf, zs...),
	}
}

func MakeMutableSetUniqueFD(zs ...*Zettel) MutableSet {
	kf := func(z *Zettel) string {
		return z.String()
	}

	return MakeMutableSet(kf, zs...)
}

func MakeMutableSetUniqueStored(zs ...*Zettel) MutableSet {
	kf := func(z *Zettel) string {
		return z.Named.Stored.Sha.String()
	}

	return MakeMutableSet(kf, zs...)
}

func MakeMutableSetUniqueAkte(zs ...*Zettel) MutableSet {
	kf := func(z *Zettel) string {
		return z.Named.Stored.Zettel.Akte.String()
	}

	return MakeMutableSet(kf, zs...)
}
