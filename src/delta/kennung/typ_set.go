package kennung

import "github.com/friedenberg/zit/src/bravo/collections"

type TypSet = collections.ValueSet[Typ, *Typ]

func MakeTypSet(ts ...Typ) TypSet {
	return collections.MakeValueSet[Typ, *Typ](
		ts...,
	)
}

type TypMutableSet = collections.MutableValueSet[Typ, *Typ]

func MakeTypMutableSet(ts ...Typ) TypMutableSet {
	return collections.MakeMutableValueSet[Typ, *Typ](
		ts...,
	)
}
