package kennung

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_ptr"
)

type (
	EtikettSet        = schnittstellen.SetPtrLike[Etikett, *Etikett]
	EtikettMutableSet = schnittstellen.MutableSetPtrLike[Etikett, *Etikett]
)

var EtikettSetEmpty EtikettSet

func init() {
	collections_ptr.RegisterGobValue[Etikett](nil)
	EtikettSetEmpty = MakeEtikettSet()
}

func MakeEtikettSet(es ...Etikett) (s EtikettSet) {
	if len(es) == 0 && EtikettSetEmpty != nil {
		return EtikettSetEmpty
	}

	return EtikettSet(
		collections_ptr.MakeValueSetValue(nil, es...),
	)
}

func MakeEtikettSetStrings(vs ...string) (s EtikettSet, err error) {
	return collections_ptr.MakeValueSetString[Etikett](nil, vs...)
}

func MakeMutableEtikettSet(hs ...Etikett) EtikettMutableSet {
	return MakeEtikettMutableSet(hs...)
}

func MakeEtikettMutableSet(hs ...Etikett) EtikettMutableSet {
	return EtikettMutableSet(
		collections_ptr.MakeMutableValueSetValue(
			nil,
			hs...,
		),
	)
}

func EtikettSetEquals(a, b EtikettSet) bool {
	return iter.SetEqualsPtr(a, b)
}
