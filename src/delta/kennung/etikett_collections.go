package kennung

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/collections2"
)

func init() {
	collections.RegisterGob[Etikett]()
}

type (
	EtikettSet        = schnittstellen.SetPtrLike[Etikett, *Etikett]
	EtikettMutableSet = schnittstellen.MutableSetPtrLike[Etikett, *Etikett]
)

func MakeEtikettSet(es ...Etikett) (s EtikettSet) {
	return EtikettSet(
		collections2.MakeValueSetValue[Etikett, *Etikett](nil, es...),
	)
}

func MakeSetStrings(vs ...string) (s EtikettSet, err error) {
	return collections2.MakeValueSetString[Etikett, *Etikett](nil, vs...)
}

func MakeMutableEtikettSet(hs ...Etikett) EtikettMutableSet {
	return MakeEtikettMutableSet(hs...)
}

func MakeEtikettMutableSet(hs ...Etikett) EtikettMutableSet {
	return EtikettMutableSet(
		collections2.MakeMutableValueSetValue[Etikett, *Etikett](
			nil,
			hs...,
		),
	)
}
