package konfig

import (
	"encoding/gob"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/hotel/etikett"
	"github.com/friedenberg/zit/src/hotel/kasten"
	"github.com/friedenberg/zit/src/hotel/typ"
)

func init() {
	gob.RegisterName("typSet", makeCompiledTypSet(nil))
	gob.RegisterName("etikettSet", makeCompiledEtikettSet(nil))
	gob.RegisterName("kastenSet", makeCompiledKastenSet(nil))
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
	return collections.MakeMutableSetPtrValueCustom[kasten.Transacted, *kasten.Transacted](
		func(k kasten.Transacted) string {
			return k.GetKennungLike().String()
		},
		s1...,
	)
}

func makeCompiledEtikettSetFromSlice(
	s1 []etikett.Transacted,
) schnittstellen.MutableSetLike[etikett.Transacted] {
	return collections.MakeMutableSetPtrValueCustom[etikett.Transacted, *etikett.Transacted](
		func(k etikett.Transacted) string {
			return k.GetKennungLike().String()
		},
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
	return collections.MakeMutableSetPtrValueCustom[typ.Transacted, *typ.Transacted](
		func(k typ.Transacted) string {
			return k.GetKennungLike().String()
		},
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
