package konfig

import (
	"encoding/gob"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections_ptr"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/transacted"
)

func init() {
	gob.RegisterName("typSet", makeCompiledTypSet(nil))
	gob.RegisterName("etikettSet", makeCompiledEtikettSet(nil))
	gob.RegisterName("kastenSet", makeCompiledKastenSet(nil))
	gob.Register(KennungKeyer[transacted.Typ, *transacted.Typ]{})
	gob.Register(KennungKeyer[transacted.Etikett, *transacted.Etikett]{})
	gob.Register(KennungKeyer[transacted.Kasten, *transacted.Kasten]{})
}

type kennungGetter interface {
	GetKennungLike() kennung.Kennung
}

type KennungKeyer[
	T kennungGetter,
	TPtr interface {
		schnittstellen.Ptr[T]
		kennungGetter
	},
] struct{}

func (sk KennungKeyer[T, TPtr]) GetKey(e T) string {
	return e.GetKennungLike().String()
}

func (sk KennungKeyer[T, TPtr]) GetKeyPtr(e TPtr) string {
	if e == nil {
		return ""
	}

	return e.GetKennungLike().String()
}

func makeCompiledKastenSet(
	s1 schnittstellen.SetLike[transacted.Kasten],
) schnittstellen.MutableSetLike[transacted.Kasten] {
	if s1 == nil {
		return makeCompiledKastenSetFromSlice(nil)
	}

	return s1.CloneMutableSetLike()
}

func makeCompiledKastenSetFromSlice(
	s1 []transacted.Kasten,
) schnittstellen.MutableSetLike[transacted.Kasten] {
	return collections_ptr.MakeMutableSetValue[transacted.Kasten, *transacted.Kasten](
		KennungKeyer[transacted.Kasten, *transacted.Kasten]{},
		s1...,
	)
}

func makeCompiledEtikettSetFromSlice(
	s1 []transacted.Etikett,
) schnittstellen.MutableSetLike[transacted.Etikett] {
	return collections_ptr.MakeMutableSetValue[transacted.Etikett, *transacted.Etikett](
		KennungKeyer[transacted.Etikett, *transacted.Etikett]{},
		s1...,
	)
}

func makeCompiledEtikettSet(
	s1 schnittstellen.SetLike[transacted.Etikett],
) schnittstellen.MutableSetLike[transacted.Etikett] {
	if s1 == nil {
		return makeCompiledEtikettSetFromSlice(nil)
	}

	return s1.CloneMutableSetLike()
}

func makeCompiledTypSetFromSlice(
	s1 []transacted.Typ,
) schnittstellen.MutableSetLike[transacted.Typ] {
	return collections_ptr.MakeMutableSetValue[transacted.Typ, *transacted.Typ](
		KennungKeyer[transacted.Typ, *transacted.Typ]{},
		s1...,
	)
}

func makeCompiledTypSet(
	s1 schnittstellen.SetLike[transacted.Typ],
) schnittstellen.MutableSetLike[transacted.Typ] {
	if s1 == nil {
		return makeCompiledTypSetFromSlice(nil)
	}

	return s1.CloneMutableSetLike()
}
