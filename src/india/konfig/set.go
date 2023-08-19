package konfig

import (
	"encoding/gob"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections_ptr"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/sku"
)

func init() {
	gob.RegisterName("typSet", makeCompiledTypSet(nil))
	gob.RegisterName("etikettSet", makeCompiledEtikettSet(nil))
	gob.RegisterName("kastenSet", makeCompiledKastenSet(nil))
	gob.Register(KennungKeyer[sku.TransactedTyp, *sku.TransactedTyp]{})
	gob.Register(KennungKeyer[sku.TransactedEtikett, *sku.TransactedEtikett]{})
	gob.Register(KennungKeyer[sku.TransactedKasten, *sku.TransactedKasten]{})
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
	s1 schnittstellen.SetLike[sku.TransactedKasten],
) schnittstellen.MutableSetLike[sku.TransactedKasten] {
	if s1 == nil {
		return makeCompiledKastenSetFromSlice(nil)
	}

	return s1.CloneMutableSetLike()
}

func makeCompiledKastenSetFromSlice(
	s1 []sku.TransactedKasten,
) schnittstellen.MutableSetLike[sku.TransactedKasten] {
	return collections_ptr.MakeMutableSetValue[sku.TransactedKasten, *sku.TransactedKasten](
		KennungKeyer[sku.TransactedKasten, *sku.TransactedKasten]{},
		s1...,
	)
}

func makeCompiledEtikettSetFromSlice(
	s1 []sku.TransactedEtikett,
) schnittstellen.MutableSetLike[sku.TransactedEtikett] {
	return collections_ptr.MakeMutableSetValue[sku.TransactedEtikett, *sku.TransactedEtikett](
		KennungKeyer[sku.TransactedEtikett, *sku.TransactedEtikett]{},
		s1...,
	)
}

func makeCompiledEtikettSet(
	s1 schnittstellen.SetLike[sku.TransactedEtikett],
) schnittstellen.MutableSetLike[sku.TransactedEtikett] {
	if s1 == nil {
		return makeCompiledEtikettSetFromSlice(nil)
	}

	return s1.CloneMutableSetLike()
}

func makeCompiledTypSetFromSlice(
	s1 []sku.TransactedTyp,
) schnittstellen.MutableSetLike[sku.TransactedTyp] {
	return collections_ptr.MakeMutableSetValue[sku.TransactedTyp, *sku.TransactedTyp](
		KennungKeyer[sku.TransactedTyp, *sku.TransactedTyp]{},
		s1...,
	)
}

func makeCompiledTypSet(
	s1 schnittstellen.SetLike[sku.TransactedTyp],
) schnittstellen.MutableSetLike[sku.TransactedTyp] {
	if s1 == nil {
		return makeCompiledTypSetFromSlice(nil)
	}

	return s1.CloneMutableSetLike()
}
