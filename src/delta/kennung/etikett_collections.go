package kennung

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
)

func init() {
	collections.RegisterGob[Etikett]()
}

type (
	EtikettSet        = schnittstellen.SetLike[Etikett]
	EtikettMutableSet = schnittstellen.MutableSetLike[Etikett]
)

func MakeEtikettSet(es ...Etikett) (s EtikettSet) {
	return EtikettSet(collections.MakeSet((Etikett).String, es...))
}

func MakeSetStrings(vs ...string) (s EtikettSet, err error) {
	f := collections.MakeFlagCommasFromExisting(
		collections.SetterPolicyAppend,
		&s,
	)

	err = f.SetMany(vs...)

	return
}

func MakeEtikettMutableSet(hs ...Etikett) EtikettMutableSet {
	return EtikettMutableSet(
		collections.MakeMutableSet[Etikett](
			(Etikett).String,
			hs...,
		),
	)
}
