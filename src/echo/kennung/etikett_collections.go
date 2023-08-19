package kennung

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections_ptr"
)

func init() {
	collections_ptr.RegisterGobValue[Etikett, *Etikett](nil)
}

type (
	EtikettSet        = schnittstellen.SetPtrLike[Etikett, *Etikett]
	EtikettMutableSet = schnittstellen.MutableSetPtrLike[Etikett, *Etikett]
)

func MakeEtikettSet(es ...Etikett) (s EtikettSet) {
	return EtikettSet(
		collections_ptr.MakeValueSetValue[Etikett, *Etikett](nil, es...),
	)
}

func MakeEtikettSetStrings(vs ...string) (s EtikettSet, err error) {
	return collections_ptr.MakeValueSetString[Etikett, *Etikett](nil, vs...)
}

func MakeMutableEtikettSet(hs ...Etikett) EtikettMutableSet {
	return MakeEtikettMutableSet(hs...)
}

func MakeEtikettMutableSet(hs ...Etikett) EtikettMutableSet {
	return EtikettMutableSet(
		collections_ptr.MakeMutableValueSetValue[Etikett, *Etikett](
			nil,
			hs...,
		),
	)
}
