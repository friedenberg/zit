package konfig

import (
	"encoding/gob"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections2"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/hotel/etikett"
	"github.com/friedenberg/zit/src/hotel/kasten"
	"github.com/friedenberg/zit/src/hotel/typ"
)

func init() {
	gob.RegisterName("typSet", makeCompiledTypSet(nil))
	gob.RegisterName("etikettSet", makeCompiledEtikettSet(nil))
	gob.RegisterName("kastenSet", makeCompiledKastenSet(nil))
	gob.Register(KennungKeyer[typ.Transacted, *typ.Transacted]{})
	gob.Register(KennungKeyer[etikett.Transacted, *etikett.Transacted]{})
	gob.Register(KennungKeyer[kasten.Transacted, *kasten.Transacted]{})
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
	s1 schnittstellen.Set[kasten.Transacted],
) schnittstellen.MutableSetLike[kasten.Transacted] {
	if s1 == nil {
		return makeCompiledKastenSetFromSlice(nil)
	}

	return s1.CloneMutableSetLike()
}

func makeCompiledKastenSetFromSlice(
	s1 []kasten.Transacted,
) schnittstellen.MutableSetLike[kasten.Transacted] {
	return collections2.MakeMutableSetValue[kasten.Transacted, *kasten.Transacted](
		KennungKeyer[kasten.Transacted, *kasten.Transacted]{},
		s1...,
	)
}

func makeCompiledEtikettSetFromSlice(
	s1 []etikett.Transacted,
) schnittstellen.MutableSetLike[etikett.Transacted] {
	return collections2.MakeMutableSetValue[etikett.Transacted, *etikett.Transacted](
		KennungKeyer[etikett.Transacted, *etikett.Transacted]{},
		s1...,
	)
}

func makeCompiledEtikettSet(
	s1 schnittstellen.Set[etikett.Transacted],
) schnittstellen.MutableSetLike[etikett.Transacted] {
	if s1 == nil {
		return makeCompiledEtikettSetFromSlice(nil)
	}

	return s1.CloneMutableSetLike()
}

func makeCompiledTypSetFromSlice(
	s1 []typ.Transacted,
) schnittstellen.MutableSetLike[typ.Transacted] {
	return collections2.MakeMutableSetValue[typ.Transacted, *typ.Transacted](
		KennungKeyer[typ.Transacted, *typ.Transacted]{},
		s1...,
	)
}

func makeCompiledTypSet(
	s1 schnittstellen.Set[typ.Transacted],
) schnittstellen.MutableSetLike[typ.Transacted] {
	if s1 == nil {
		return makeCompiledTypSetFromSlice(nil)
	}

	return s1.CloneMutableSetLike()
}
