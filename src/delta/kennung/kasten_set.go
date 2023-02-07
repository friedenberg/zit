package kennung

import "github.com/friedenberg/zit/src/charlie/collections"

type KastenSet = collections.ValueSet[Kasten, *Kasten]

func MakeKastenSet(ts ...Kasten) KastenSet {
	return collections.MakeValueSet[Kasten, *Kasten](
		ts...,
	)
}

type KastenMutableSet = collections.MutableValueSet[Kasten, *Kasten]

func MakeKastenMutableSet(ts ...Kasten) KastenMutableSet {
	return collections.MakeMutableValueSet[Kasten, *Kasten](
		ts...,
	)
}
