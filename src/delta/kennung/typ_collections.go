package kennung

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/collections2"
)

func init() {
	collections.RegisterGob[Typ]()
}

type (
	TypSet        = schnittstellen.SetPtrLike[Typ, *Typ]
	TypMutableSet = schnittstellen.MutableSetPtrLike[Typ, *Typ]
)

func MakeTypSet(es ...Typ) (s TypSet) {
	return TypSet(
		collections2.MakeValueSetValue[Typ, *Typ](nil, es...),
	)
}

func MakeTypSetStrings(vs ...string) (s TypSet, err error) {
	return collections2.MakeValueSetString[Typ, *Typ](nil, vs...)
}

func MakeMutableTypSet(hs ...Typ) TypMutableSet {
	return MakeTypMutableSet(hs...)
}

func MakeTypMutableSet(hs ...Typ) TypMutableSet {
	return TypMutableSet(
		collections2.MakeMutableValueSetValue[Typ, *Typ](
			nil,
			hs...,
		),
	)
}
