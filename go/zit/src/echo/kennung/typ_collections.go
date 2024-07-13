package kennung

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_ptr"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/echo/thyme"
)

func init() {
	collections_value.RegisterGobValue[thyme.Time](nil)
}

type (
	TypSet        = interfaces.SetPtrLike[Typ, *Typ]
	TypMutableSet = interfaces.MutableSetPtrLike[Typ, *Typ]
)

func MakeTypSet(es ...Typ) (s TypSet) {
	return TypSet(
		collections_ptr.MakeValueSetValue[Typ, *Typ](nil, es...),
	)
}

func MakeTypSetStrings(vs ...string) (s TypSet, err error) {
	return collections_ptr.MakeValueSetString[Typ, *Typ](nil, vs...)
}

func MakeMutableTypSet(hs ...Typ) TypMutableSet {
	return MakeTypMutableSet(hs...)
}

func MakeTypMutableSet(hs ...Typ) TypMutableSet {
	return TypMutableSet(
		collections_ptr.MakeMutableValueSetValue[Typ, *Typ](
			nil,
			hs...,
		),
	)
}
